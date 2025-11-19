package email

import (
	"context"
	"encoding/json"
	"errors"

	"go.etcd.io/bbolt"
)

var bucketName = []byte("reset_codes")

// BoltRepository implementa Repository con BoltDB.
type BoltRepository struct {
	db *bbolt.DB
}

// initBucketIfMissing crea el bucket si no existe.
func initBucketIfMissing(db *bbolt.DB) error {
	return db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucketName)
		return err
	})
}

func (r *BoltRepository) SaveCode(ctx context.Context, entry CodeEntry) error {
	// usamos Update (transacción)
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

func (r *BoltRepository) GetCodeEntry(ctx context.Context, email string) (*CodeEntry, error) {
	var entry CodeEntry
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
