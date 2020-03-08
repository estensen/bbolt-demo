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
		s.testDB()
		w.WriteHeader(http.StatusOK)
	}
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
