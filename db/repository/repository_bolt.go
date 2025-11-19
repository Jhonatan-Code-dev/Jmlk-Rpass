package repository

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/Jhonatan-Code-dev/Jmlk-Rpass/db/models"
	"go.etcd.io/bbolt"
)

var bucketName = []byte("reset_codes")

type BoltRepository struct {
	db *bbolt.DB
}

func InitBucketIfMissing(db *bbolt.DB) error {
	return db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucketName)
		return err
	})
}

func (r *BoltRepository) SaveCode(ctx context.Context, entry models.CodeEntry) error {
	return r.db.Update(func(tx *bbolt.Tx) error {
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
	err := r.db.View(func(tx *bbolt.Tx) error {
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
	return r.db.Close()
}
