package repository

import (
	"log"
	"os"
	"path/filepath"

	"github.com/Jhonatan-Code-dev/Jmlk-Rpass/config"
	"go.etcd.io/bbolt"
)

func InitBoltRepository(cfg *config.EmailConfig) *BoltRepository {
	baseDir, _ := os.Getwd()
	dbPath := filepath.Join(baseDir, cfg.DatabaseFolder, cfg.DatabaseName)

	db := mustOpenBoltDB(dbPath)
	mustInitBucket(db)

	return &BoltRepository{DB: db}
}

func mustOpenBoltDB(dbPath string) *bbolt.DB {
	if err := os.MkdirAll(filepath.Dir(dbPath), os.ModePerm); err != nil {
		log.Fatalf("mkdir storage: %v", err)
	}

	db, err := bbolt.Open(dbPath, 0666, nil)
	if err != nil {
		log.Fatalf("open bolt db: %v", err)
	}

	return db
}

func mustInitBucket(db *bbolt.DB) {
	err := db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("reset_codes"))
		return err
	})

	if err != nil {
		log.Fatalf("init bucket: %v", err)
	}
}
