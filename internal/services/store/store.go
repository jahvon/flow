package store

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
	bolt "go.etcd.io/bbolt"

	"github.com/jahvon/flow/internal/filesystem"
)

const (
	BucketEnv  = "FLOW_PROCESS_BUCKET"
	RootBucket = "root"

	storeFileName = "store.db"
)

type Store struct {
	db *bolt.DB
}

func NewStore() (*Store, error) {
	db, err := bolt.Open(Path(), 0666, &bolt.Options{Timeout: 5 * time.Second})
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}
	return &Store{db: db}, nil
}

// CreateBucket creates the current process bucket if it doesn't exist
func (s *Store) CreateBucket() error {
	return s.db.Update(func(tx *bolt.Tx) error {
		id := bucketID()
		_, err := tx.CreateBucketIfNotExists(id)
		if err != nil {
			return fmt.Errorf("failed to create bucket %s: %w", id, err)
		}
		return nil
	})
}

// DeleteBucket deletes the current process bucket
func (s *Store) DeleteBucket() error {
	return s.db.Update(func(tx *bolt.Tx) error {
		id := bucketID()
		err := tx.DeleteBucket(id)
		if err != nil {
			if errors.Is(err, bolt.ErrBucketNotFound) {
				return nil
			}
			return fmt.Errorf("failed to delete bucket %s: %w", id, err)
		}
		return nil
	})
}

// Set stores a key-value pair in the process bucket
func (s *Store) Set(key, value string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		id := bucketID()
		bucket := tx.Bucket(id)
		if bucket == nil {
			return fmt.Errorf("bucket %s not found", id)
		}
		err := bucket.Put([]byte(key), []byte(value))
		if err != nil {
			return fmt.Errorf("failed to put key-value pair for key %s in bucket %s: %w", key, id, err)
		}
		return nil
	})
}

// Get retrieves a value for a key from the process bucket
func (s *Store) Get(key string) (string, error) {
	var value []byte
	err := s.db.View(func(tx *bolt.Tx) error {
		id := bucketID()
		bucket := tx.Bucket(id)
		if bucket == nil {
			return fmt.Errorf("bucket %s not found", id)
		}
		value = bucket.Get([]byte(key))
		if value == nil {
			rBucket := tx.Bucket([]byte(RootBucket))
			if rBucket != nil {
				value = rBucket.Get([]byte(key))
			}
		}
		if value == nil {
			return fmt.Errorf("key %s not found in bucket %s", key, id)
		}
		return nil
	})
	return string(value), err
}

// Delete removes a key from the process bucket
func (s *Store) Delete(key string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		id := bucketID()
		bucket := tx.Bucket(id)
		if bucket == nil {
			return fmt.Errorf("bucket %s not found", id)
		}
		err := bucket.Delete([]byte(key))
		if err != nil {
			return fmt.Errorf("failed to delete key %s from bucket %s: %w", key, id, err)
		}
		return nil
	})
}

func (s *Store) Close() error {
	return s.db.Close()
}

func Path() string {
	cacheDir := filesystem.CachedDataDirPath()
	return filepath.Join(cacheDir, storeFileName)
}

func SetProcessBucketID(id string, force bool) error {
	if _, set := os.LookupEnv(BucketEnv); set && !force {
		return nil
	}

	replacer := strings.NewReplacer(":", "_", "/", "_", " ", "_")
	id = replacer.Replace(id)
	return os.Setenv(BucketEnv, id)
}

func bucketID() []byte {
	processBucket, set := os.LookupEnv(BucketEnv)
	if !set {
		return []byte(RootBucket)
	}
	return []byte(processBucket)
}