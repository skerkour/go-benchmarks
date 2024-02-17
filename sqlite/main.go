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
	CreatedAt Time
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
	os.RemoveAll("./test.db")
	os.RemoveAll("./test.db-shm")
	os.RemoveAll("./test.db-wal")

	uuid.EnableRandPool()

	connectionUrlParams := make(url.Values)
	connectionUrlParams.Add("_txlock", "immediate")
	connectionUrlParams.Add("_journal_mode", "WAL")
	connectionUrlParams.Add("_busy_timeout", "5000")
	connectionUrlParams.Add("_synchronous", "NORMAL")
	connectionUrlParams.Add("_cache_size", "1000000000")
	connectionUrlParams.Add("_foreign_keys", "true")
	connectionUrl, err := url.Parse("./test.db?" + connectionUrlParams.Encode())
	if err != nil {
		log.Fatal(err)
	}

	writeDB, err := sql.Open("sqlite3", connectionUrl.String())
	if err != nil {
		log.Fatal(err)
	}
	defer writeDB.Close()
	writeDB.SetMaxOpenConns(1)
	err = setupSqlite(writeDB)
	if err != nil {
		log.Fatal(err)
	}

	readDB, err := sql.Open("sqlite3", connectionUrl.String())
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
		created_at INTEGER NOT NULL,
		counter INT NOT NULL
	) STRICT`)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("starting")

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
			// we use a goroutine-local counter to avoid the performance impact of updating an atomic counter
			var readsLocal int64

			for {
				var something entity

				if len(ticker.C) > 0 {
					break
				}

				rows, err := readDB.Query("SELECT * FROM test LIMIT 1")
				if err != nil {
					log.Fatal(err)
				}
				if !rows.Next() {
					rows.Close()
					continue
				}

				rows.Scan(&something.ID, &something.CreatedAt, &something.Counter)
				rows.Close()
				readsLocal += 1
			}
			reads.Add(readsLocal)
			wg.Done()
		}()
	}

	wg.Add(concurrentWriters)
	for c := 0; c < concurrentWriters; c += 1 {
		go func() {
			startUnixMilliseconds := start.UnixMilli()
			// we use a goroutine-local counter to avoid the performance impact of updating an atomic counter
			var writesLocal int64

			for {
				if len(ticker.C) > 0 {
					break
				}

				uu := uuid.New()

				_, err = writeDB.Exec(`INSERT INTO test
					(id, created_at, counter) VALUES (?, ?, ?)`, uu[:], startUnixMilliseconds, writes.Load())
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

	log.Println("elapsed:", elapsed)
	fmt.Println("----------------------")

	log.Printf("%d reads\n", reads.Load())

	throughputRead := float64(reads.Load()) / float64(elapsed.Seconds())
	log.Printf("%f reads/s\n", throughputRead)

	fmt.Println("----------------------")

	log.Printf("%d writes\n", writes.Load())
	throughputWrite := float64(writes.Load()) / float64(elapsed.Seconds())
	log.Printf("%f writes/s\n", throughputWrite)
}
