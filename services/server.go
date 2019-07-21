package services

import (
	"encoding/json"
	"fmt"
	"gate/models"
	"net/http"
	"sync"
	"time"
)

type Server struct {
}

func NewServer(wg *sync.WaitGroup) Iservice {
	return &Server{}
}

func (s *Server) Start() {
	server := &http.Server{
		Addr:         ":8000",
		ReadTimeout:  10 * time.Minute,
		WriteTimeout: 10 * time.Minute,
	}
	http.HandleFunc("/", IndexHandler)
	server.ListenAndServe()
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	indexDBs := new(models.IndexDBs)
	if err := indexDBs.Select(); err != nil {
		fmt.Printf("select indexdb error: %s\n", err.Error())
	}
	indexDBs.ToList()
	bs, err := json.Marshal(indexDBs)
	if err != nil {
		fmt.Println(err)
	}
	w.Write(bs)
}

func (s *Server) Stop() {}
