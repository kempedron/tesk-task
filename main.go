package main

import (
	"net/http"
	"strconv"
	"sync"
	Q "testTask/Queue"
	"time"
)

var (
	QMap = make(map[string]*Q.Queue)
	mu   sync.RWMutex
)

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("GET /{queue_name}", handlerGet)
	mux.HandleFunc("PUT /{queue_name}", handlerAdd)

	http.ListenAndServe("127.0.0.1:8080", mux)

}

func handlerAdd(w http.ResponseWriter, r *http.Request) {
	qName := r.PathValue("queue_name")
	qText := r.URL.Query().Get("v")

	mu.Lock()
	q, ok := QMap[qName]
	if !ok {
		q = Q.InitNewQ(qName)
		QMap[qName] = q
	}
	mu.Unlock()
	q.AddQ(qText)

	w.WriteHeader(http.StatusOK)
}

func handlerGet(w http.ResponseWriter, r *http.Request) {

	timeoutStr := r.URL.Query().Get("timeout")
	seconds, _ := strconv.Atoi(timeoutStr)
	timeout := time.Duration(seconds) * time.Second

	qName := r.PathValue("queue_name")
	mu.RLock()

	q, ok := QMap[qName]
	mu.RUnlock()
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	tryGet := func() bool {
		if node := q.GetQ(); node != nil {
			w.Write([]byte(node.Val))
			return true
		}
		return false
	}

	if tryGet() {
		return
	}

	if timeout <= 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	timer := time.NewTimer(timeout)

	defer timer.Stop()

	for {
		select {
		case <-q.NotifyChan:
			if node := q.GetQ(); node != nil {
				w.Write([]byte(node.Val))
				return
			}
			continue

		case <-timer.C:
			w.WriteHeader(http.StatusNotFound)
		}
	}
}
