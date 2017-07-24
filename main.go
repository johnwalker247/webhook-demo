package main

import (
	// "context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

// Version - version to show
var Version = "0.0.14"

// Port - default port to start application on
var Port = ":8090"

// ReceivedWebhook - keeps info about received webhook
type ReceivedWebhook struct {
	Payload    string
	ReceivedAt time.Time
}

func main() {
	// preparing HTTP server
	srv := &http.Server{Addr: Port, Handler: http.DefaultServeMux}

	// all received webhooks will be stored in this variable
	var webhooks []*ReceivedWebhook
	mu := &sync.RWMutex{}

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// handler to display all received webhooks
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Version: %s \n", Version)
		fmt.Fprintln(w, "Received webhooks:")
		mu.RLock()
		for _, v := range webhooks {
			fmt.Fprintln(w, fmt.Sprintf("%s: %s", v.ReceivedAt.String(), v.Payload))
		}
		mu.RUnlock()
	})

	// incomming webhook handler
	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		bd, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		defer r.Body.Close()

		// appending webhook to received list
		mu.Lock()
		webhooks = append(webhooks, &ReceivedWebhook{
			ReceivedAt: time.Now(),
			Payload:    string(bd),
		})
		mu.Unlock()
		w.WriteHeader(http.StatusOK)
	})

	fmt.Println("Version: ", Version)
	fmt.Println("Listening on ", Port)
	// starting server
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}
