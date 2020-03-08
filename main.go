package main

import (
	"fmt"
	"log"
	bolt "go.etcd.io/bbolt"

	"github.com/estensen/bbolt-demo/db"
)

type server struct {
	db *db.DB
}

func main() {
	db, err := db.NewDB("answers.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	srv := server{
		db: db,
	}

	srv.testDB()
}

func (s *server) testDB() {
	_ = s.putAnswer("42")
	answer, _ := s.getAnswer()
	fmt.Println(answer)
}

func (s *server) putAnswer(answer string) (err error) {
	err = s.db.Update(func(tx *bolt.Tx) error {
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
		return err
	}
	return nil
}

func (s *server) getAnswer() (answer string, err error) {
	var val []byte
	err = s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("bucket"))
		v := b.Get([]byte("answer"))
		val = append([]byte(nil), v...)

		return nil
	})
	if err != nil {
		return "", err
	}

	return string(val), nil
}

