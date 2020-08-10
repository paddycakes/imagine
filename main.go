package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/paddycakes/imagine/imagemagick"
	"github.com/paddycakes/imagine/model"
)

func main() {
	http.HandleFunc("/", HandlePubSub)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	// Start HTTP server.
	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

// Receives and processes a Pub/Sub push message.
func HandlePubSub(w http.ResponseWriter, r *http.Request) {
	var m model.PubSubMessage
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("ioutil.ReadAll: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(body, &m); err != nil {
		log.Printf("json.Unmarshal: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	var e model.GCSEvent
	if err := json.Unmarshal(m.Message.Data, &e); err != nil {
		log.Printf("json.Unmarshal: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if e.Name == "" || e.Bucket == "" {
		log.Printf("invalid GCSEvent: expected name and bucket")
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if err := imagemagick.BlurOffensiveImages(r.Context(), e); err != nil {
		log.Printf("imagemagick.BlurOffensiveImages: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
