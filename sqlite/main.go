package main

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"log"
	"net/url"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

// Time is used to store timestamps as INT in SQLite
type Time int64

func (t *Time) Scan(val any) (err error) {
	switch v := val.(type) {
	case int64:
		*t = Time(v)
		return nil
	case string:
		tt, err := time.Parse(time.RFC3339, v)
		if err != nil {
			return err
		}
		*t = Time(tt.UnixMilli())
		return nil
	default:
		return fmt.Errorf("Time.Scan: Unsupported type: %T", v)
	}
}

func (t *Time) Value() (driver.Value, error) {
	return *t, nil
}

type entity struct {
	ID        uuid.UUID
	Timestamp Time
	Counter   int64
}

func setupSqlite(db *sql.DB) (err error) {
	pragmas := []string{
		// "journal_mode = WAL",
		// "busy_timeout = 5000",
		// "synchronous = NORMAL",
		// "cache_size = 1000000000", // 1GB
		// "foreign_keys = true",
		"temp_store = memory",
		// "mmap_size = 3000000000",
	}

	for _, pragma := range pragmas {
		_, err = db.Exec("PRAGMA " + pragma)
		if err != nil {
			return
		}
	}

	return nil
}

func main() {
	cleanup()
	defer cleanup()

	uuid.EnableRandPool()

	connectionUrlParams := make(url.Values)
	connectionUrlParams.Add("_txlock", "immediate")
	connectionUrlParams.Add("_journal_mode", "WAL")
	connectionUrlParams.Add("_busy_timeout", "5000")
	connectionUrlParams.Add("_synchronous", "NORMAL")
	connectionUrlParams.Add("_cache_size", "1000000000")
	connectionUrlParams.Add("_foreign_keys", "true")
	connectionUrl := "file:test.db?" + connectionUrlParams.Encode()

	writeDB, err := sql.Open("sqlite3", connectionUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer writeDB.Close()
	writeDB.SetMaxOpenConns(1)
	err = setupSqlite(writeDB)
	if err != nil {
		log.Fatal(err)
	}

	readDB, err := sql.Open("sqlite3", connectionUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer readDB.Close()
	readDB.SetMaxOpenConns(max(4, runtime.NumCPU()))
	err = setupSqlite(readDB)
	if err != nil {
		log.Fatal(err)
	}

	_, err = writeDB.Exec(`CREATE TABLE test (
		id BLOB NOT NULL PRIMARY KEY,
		timestamp INTEGER NOT NULL,
		counter INT NOT NULL
	) STRICT`)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Inserting 5,000,000 rows")
	err = setupDB(writeDB)
	if err != nil {
		log.Fatal(err)
	}

	var recordIdToFind uuid.UUID
	row := readDB.QueryRow("SELECT id FROM test ORDER BY id DESC LIMIT 1")
	if row.Err() != nil {
		log.Fatal(row.Err())
	}
	row.Scan(&recordIdToFind)

	log.Println("Starting benchmark")

	concurrentReaders := 500
	concurrentWriters := 1
	var wg sync.WaitGroup
	var reads atomic.Int64
	var writes atomic.Int64
	ticker := time.NewTicker(10 * time.Second)
	start := time.Now()

	wg.Add(concurrentReaders)
	for c := 0; c < concurrentReaders; c += 1 {
		go func() {
			// we use a goroutine-local counter to avoid the performance impact of updating a shared atomic counter
			var readsLocal int64

			for {
				var record entity

				if len(ticker.C) > 0 {
					break
				}

				row := readDB.QueryRow("SELECT * FROM test WHERE id = ?", recordIdToFind)
				if row.Err() != nil {
					log.Fatal(row.Err())
				}

				row.Scan(&record.ID, &record.Timestamp, &record.Counter)
				readsLocal += 1
			}
			reads.Add(readsLocal)
			wg.Done()
		}()
	}

	wg.Add(concurrentWriters)
	for c := 0; c < concurrentWriters; c += 1 {
		go func() {
			timestamp := start.UnixMilli()
			// we use a goroutine-local counter to avoid the performance impact of updating a shared atomic counter
			var writesLocal int64

			for {
				if len(ticker.C) > 0 {
					break
				}

				recordID := uuid.Must(uuid.NewV7())

				_, err = writeDB.Exec(`INSERT INTO test
					(id, timestamp, counter) VALUES (?, ?, ?)`, recordID[:], timestamp, writesLocal)
				if err != nil {
					log.Fatal(err)
				}
				writesLocal += 1
			}
			writes.Add(writesLocal)
			wg.Done()
		}()
	}

	wg.Wait()

	elapsed := time.Since(start)

	log.Println("Benchmark stopped:", elapsed)
	fmt.Println("----------------------")

	log.Printf("%d reads\n", reads.Load())

	throughputRead := float64(reads.Load()) / float64(elapsed.Seconds())
	log.Printf("%f reads/s\n", throughputRead)

	fmt.Println("----------------------")

	log.Printf("%d writes\n", writes.Load())
	throughputWrite := float64(writes.Load()) / float64(elapsed.Seconds())
	log.Printf("%f writes/s\n", throughputWrite)
}

func cleanup() {
	os.RemoveAll("./test.db")
	os.RemoveAll("./test.db-shm")
	os.RemoveAll("./test.db-wal")
}

func setupDB(db *sql.DB) (err error) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	timestamp := time.Now().UTC().UnixMilli()
	for i := range 5_000_000 {
		recordID := uuid.Must(uuid.NewV7())
		_, err = tx.Exec(`INSERT INTO test
					(id, timestamp, counter) VALUES (?, ?, ?)`, recordID[:], timestamp, i)
		if err != nil {
			return err
		}

		// insert by batches of 500,000 rows
		if i%500_000 == 0 {
			err = tx.Commit()
			if err != nil {
				tx.Rollback()
				return err
			}
			tx, err = db.Begin()
			if err != nil {
				return err
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	_, err = db.Exec("VACUUM")
	if err != nil {
		return err
	}

	_, err = db.Exec("ANALYZE")
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
