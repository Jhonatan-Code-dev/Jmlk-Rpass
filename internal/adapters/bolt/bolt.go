// internal/adapters/bolt/bolt.go
package bolt

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Jhonatan-Code-dev/Jmlk-Rpass/internal/domain"
	"github.com/Jhonatan-Code-dev/Jmlk-Rpass/internal/ports"
	"go.etcd.io/bbolt"
)

var bucketName = []byte("reset_codes")

type BoltRepo struct {
	db *bbolt.DB
}

func New(path string) (ports.Repository, error) {
	if path == "" {
		return nil, fmt.Errorf("path required")
	}
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("mkdir: %w", err)
	}
	db, err := bbolt.Open(path, 0666, nil)
	if err != nil {
		return nil, fmt.Errorf("open bolt: %w", err)
	}
	if err := db.Update(func(tx *bbolt.Tx) error {
		_, e := tx.CreateBucketIfNotExists(bucketName)
		return e
	}); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("init bucket: %w", err)
	}
	return &BoltRepo{db: db}, nil
}

func (r *BoltRepo) SaveCode(ctx context.Context, entry domain.CodeEntry) error {
	return r.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucketName)
		if b == nil {
			return fmt.Errorf("bucket missing")
		}
		data, err := json.Marshal(entry)
		if err != nil {
			return err
		}
		return b.Put([]byte(entry.Email), data)
	})
}

func (r *BoltRepo) GetCodeEntry(ctx context.Context, email string) (*domain.CodeEntry, error) {
	var e domain.CodeEntry
	err := r.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucketName)
		if b == nil {
			return nil
		}
		raw := b.Get([]byte(email))
		if raw == nil {
			return nil
		}
		return json.Unmarshal(raw, &e)
	})
	if err != nil {
		return nil, err
	}
	if e.Email == "" {
		return nil, nil
	}
	return &e, nil
}

func (r *BoltRepo) Close() error {
	return r.db.Close()
}
