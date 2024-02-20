// Benchmark the difference of performance between INT and TEXT enums with SQLite
package sqliteenum

import (
	"database/sql"
	"log"
	"net/url"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func BenchmarkSqliteEnum(b *testing.B) {
	cleanup()
	defer cleanup()

	db, err := setupSqlite(":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	b.Run("INT", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var count int64
			row := db.QueryRow("SELECT COUNT(*) FROM test_enum WHERE value_int = ?", MyEnumCompleted)
			if err != nil {
				log.Fatal(err)
			}
			row.Scan(&count)
		}
	})

	// SELECT COUNT(*) FROM (SELECT * FROM test_enum WHERE value_int = 1 LIMIT 1)
	b.Run("INT LIMIT", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var count int64
			row := db.QueryRow("SELECT COUNT(*) FROM (SELECT * FROM test_enum WHERE value_int = ? LIMIT 1)", MyEnumCompleted)
			if err != nil {
				log.Fatal(err)
			}
			row.Scan(&count)
		}
	})

	b.Run("TEXT", func(b *testing.B) {
		completedText := MyEnumCompleted.String()
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var count int64
			row := db.QueryRow("SELECT COUNT(*) FROM test_enum WHERE value_text = ?", completedText)
			if err != nil {
				log.Fatal(err)
			}
			row.Scan(&count)
		}
	})
}

type MyEnum int64

const (
	MyEnumActive MyEnum = iota
	MyEnumPending
	MyEnumProcessing
	MyEnumCanceled
	MyEnumCompleted
	MyEnumAnotherValue
	MyEnumSomethingElse
	MyEnumDraft
	MyEnumUnknown
)

func (e MyEnum) String() string {
	switch e {
	case MyEnumActive:
		return "active"
	case MyEnumPending:
		return "pending"
	case MyEnumProcessing:
		return "processing"
	case MyEnumCanceled:
		return "canceled"
	case MyEnumCompleted:
		return "completed"
	case MyEnumAnotherValue:
		return "anothervalue"
	case MyEnumSomethingElse:
		return "somethingelse"
	case MyEnumDraft:
		return "draft"
	case MyEnumUnknown:
		return "unknown"
	default:
		return "error"
	}
}

var enumValues = []MyEnum{
	MyEnumActive,
	MyEnumPending,
	MyEnumProcessing,
	MyEnumCanceled,
	MyEnumCompleted,
	MyEnumAnotherValue,
	MyEnumSomethingElse,
	MyEnumDraft,
	MyEnumUnknown,
}

func cleanup() {
	os.RemoveAll("./test.db")
	os.RemoveAll("./test.db-shm")
	os.RemoveAll("./test.db-wal")
}

func setupSqlite(dbName string) (db *sql.DB, err error) {
	connectionUrlParams := make(url.Values)
	connectionUrlParams.Add("_txlock", "immediate")
	connectionUrlParams.Add("_journal_mode", "WAL")
	connectionUrlParams.Add("_busy_timeout", "5000")
	connectionUrlParams.Add("_synchronous", "NORMAL")
	connectionUrlParams.Add("_cache_size", "1000000000")
	connectionUrlParams.Add("_foreign_keys", "true")
	connectionUrl := "file:" + dbName + "?" + connectionUrlParams.Encode()

	db, err = sql.Open("sqlite3", connectionUrl)
	if err != nil {
		return
	}

	_, err = db.Exec(`CREATE TABLE test_enum (
		value_int INT NOT NULL,
		value_text TEXT NOT NULL
	) STRICT`)
	if err != nil {
		return
	}

	tx, err := db.Begin()
	if err != nil {
		return
	}

	numberOfEnumValues := len(enumValues)
	for i := 0; i < 2_000_000; i += 1 {
		value := enumValues[i%numberOfEnumValues]

		_, err = tx.Exec("INSERT INTO test_enum (value_int, value_text) VALUES (?, ?)", value, value.String())
		if err != nil {
			return
		}

		// we insert by batches of 200_000
		if i%200_000 == 0 {
			err = tx.Commit()
			if err != nil {
				return
			}

			tx, err = db.Begin()
			if err != nil {
				return
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		return
	}

	_, err = db.Exec("CREATE INDEX index_test_enum_on_value_int ON test_enum (value_int)")
	if err != nil {
		return
	}

	_, err = db.Exec("CREATE INDEX index_test_enum_on_value_text ON test_enum (value_text)")
	if err != nil {
		return
	}

	_, err = db.Exec("ANALYZE")
	if err != nil {
		return
	}

	return
}
