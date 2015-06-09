package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
)

type DB struct {
	Handle *bolt.DB
}

// Open the BoltDB file
func (db *DB) Open(dbfile string) {
	var err error

	db.Handle, err = bolt.Open(dbfile, 0600, nil)
	if err != nil {
		log.Fatalln(err)
	}

	return
}

// Close the BoltDB file
func (db *DB) Close() {
	db.Handle.Close()
	return
}

// Store a key/value in a BoltDB bucket
func (d *DB) Store(bucket, key, value string) error {

	log.Println("Storing:", key)

	err := d.Handle.Update(func(tx *bolt.Tx) error {
		bkt, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return fmt.Errorf("Could not create bucket %q", bucket)
		}

		err = bkt.Put([]byte(key), []byte(value))
		if err != nil {
			return fmt.Errorf("Could not write key %q to bucket %q: %v", key, bucket, err)
		}

		return nil
	})

	return err
}

// Delete a key from a BoltDB bucket
func (d *DB) Delete(bucket, key string) error {

	log.Println("Deleting", key, "from bucket", bucket)

	err := d.Handle.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte(bucket))
		if bkt == nil {
			return fmt.Errorf("Could not locate bucket to delete: %q", bucket)
		}

		err := bkt.Delete([]byte(key))
		if err != nil {
			return fmt.Errorf("Could not delete bucket %q: %v", bucket, err)
		}

		return nil
	})

	return err
}

// Fetch a key from a BoltDB bucket
func (d *DB) Fetch(bucket, key string) (string, error) {

	var val string

	log.Println("Fetching", key, "from bucket", bucket)

	err := d.Handle.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte(bucket))
		if bkt == nil {
			return fmt.Errorf("Bucket %q not found!", bucket)
		}

		val = string(bkt.Get([]byte(key)))

		return nil
	})

	if val == "" {
		return "", fmt.Errorf("Key %q in bucket %q not found", key, bucket)
	}

	return val, err
}

// Fetch every key from a BoltDB bucket
func (d *DB) FetchAll(bucket string) ([]string, error) {
	var vals []string

	log.Println("Fetching all from bucket", bucket)

	err := d.Handle.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte(bucket))
		if bkt == nil {
			return fmt.Errorf("Bucket %q not found!", bucket)
		}

		bkt.ForEach(func(k, v []byte) error {
			vals = append(vals, string(v))
			return nil
		})

		if len(vals) == 0 {
			return fmt.Errorf("There are no keys in bucket", bucket)

		}

		return nil
	})

	return vals, err

}
