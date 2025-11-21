package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Jhonatan-Code-dev/Jmlk-Rpass/db/models"
	"go.etcd.io/bbolt"
)

var bucketName = []byte("reset_codes")

type BoltRepository struct {
	DB *bbolt.DB
}

func InitBoltDBPath(dbPath string) (*bbolt.DB, error) {
	if err := os.MkdirAll(filepath.Dir(dbPath), os.ModePerm); err != nil {
		return nil, fmt.Errorf("mkdir storage: %w", err)
	}

	db, err := bbolt.Open(dbPath, 0666, nil)
	if err != nil {
		return nil, fmt.Errorf("open bolt db: %w", err)
	}

	return db, nil
}

func InitBucketIfMissing(db *bbolt.DB) error {
	return db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucketName)
		return err
	})
}

func (r *BoltRepository) SaveCode(ctx context.Context, entry models.CodeEntry) error {
	return r.DB.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucketName)
		data, err := json.Marshal(entry)
		if err != nil {
			return err
		}
		return b.Put([]byte(entry.Email), data)
	})
}

var ErrNotFound = errors.New("registro no encontrado")

func (r *BoltRepository) GetCodeEntry(ctx context.Context, email string) (*models.CodeEntry, error) {
	var entry models.CodeEntry

	err := r.DB.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucketName)
		if b == nil {
			return ErrNotFound
		}

		data := b.Get([]byte(email))
		if data == nil {
			return ErrNotFound
		}

		return json.Unmarshal(data, &entry)
	})

	if err != nil {
		return nil, err
	}

	return &entry, nil
}

func (r *BoltRepository) Close() error {
	return r.DB.Close()
}
