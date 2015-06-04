package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
)

type DB struct {
	Handle *bolt.DB
}

func (db *DB) Open(dbfile string) {
	var err error

	db.Handle, err = bolt.Open(dbfile, 0600, nil)
	if err != nil {
		log.Fatalln(err)
	}

	return
}

func (db *DB) Close() {
	db.Handle.Close()
	return
}

func (d *DB) Store(bucket, key, value string) error {

	log.Println("Storing:", key)

	err := d.Handle.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return fmt.Errorf("Could not create bucket %q", bucket)
		}

		err = bucket.Put([]byte(key), []byte(value))
		if err != nil {
			return fmt.Errorf("Could not write key %q to bucket %q", key, bucket)
		}

		return nil
	})

	return err
}

func (d *DB) Fetch(bucket, key string) (string, error) {

	var val string

	log.Println("Fetching:", key)

	err := d.Handle.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucket))
		if bucket == nil {
			return fmt.Errorf("Bucket %q not found!", bucket)
		}

		val = string(bucket.Get([]byte(key)))

		return nil
	})

	if val == "" {
		return "", fmt.Errorf("Key %q in bucket %q not found", key, bucket)
	}

	return val, err
}
