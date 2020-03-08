package main

import (
	"fmt"
	"log"
	"net/http"

	bolt "go.etcd.io/bbolt"

	"github.com/estensen/bbolt-demo/db"
	"github.com/gorilla/mux"
)

type server struct {
	db     *db.DB
	router *mux.Router
}

func main() {
	db, err := db.NewDB("answers.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := mux.NewRouter()

	srv := server{
		db:     db,
		router: router,
	}

	srv.routes()

	log.Println("Server is running on port 3000")
	http.ListenAndServe(":3000", srv.router)
}

func (s *server) ServceHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) handlerTestDB() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		answer := s.testDB()
		w.WriteHeader(http.StatusOK)
		w.Write(answer)
	}
}

func (s *server) testDB() (answer []byte) {
	_ = s.putAnswer("42")
	v, _ := s.getAnswer()
	return v
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

func (s *server) getAnswer() (answer []byte, err error) {
	var val []byte
	err = s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("bucket"))
		v := b.Get([]byte("answer"))
		val = append([]byte(nil), v...)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return val, nil
}
