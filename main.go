package main

import (
	"fmt"
	"log"

	bolt "go.etcd.io/bbolt"
)

func main() {
	db, err := bolt.Open("answers.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("bucket"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		err = b.Put([]byte("answer"), []byte("42"))
		if err != nil {
			return fmt.Errorf("put answer: %s", err)
		}
		return nil
	})
	if err != nil {
		log.Print(err)
	}

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("bucket"))
		v := b.Get([]byte("answer"))
		fmt.Printf("The answer is %s\n", v)
		return nil
	})
	if err != nil {
		log.Print(err)
	}
}

