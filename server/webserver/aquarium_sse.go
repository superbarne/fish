package webserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (ws *WebServer) sseAquarium(w http.ResponseWriter, r *http.Request) {
	// validate id
	aquariumID, err := uuid.Parse(chi.URLParam(r, "aquariumID"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 not found"))
		return
	}

	// find aquarium
	_, err = ws.storage.Aquarium(aquariumID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 not found"))
		return
	}

	ctx := r.Context()
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	flusher, ok := w.(http.Flusher)
	if !ok {
		return
	}

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	fmt.Fprintf(w, "event: ping\ndata: {}\n\n")
	flusher.Flush()

	// send old fishes
	fishes, err := ws.storage.Fishes(aquariumID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 internal server error"))
		return
	}

	for _, fish := range fishes {
		raw, _ := json.Marshal(fish)
		fmt.Fprintf(w, "event: fishjoin\ndata: %s\n\n", raw)
	}
	flusher.Flush()

	// subscribe to new fishes
	newFishes := ws.pubsub.Subscribe("aquarium:"+aquariumID.String(), ctx, 10)
	defer ws.pubsub.Unsubscribe("aquarium:"+aquariumID.String(), ctx)
	deleteFishes := ws.pubsub.Subscribe("aquarium:"+aquariumID.String()+":delete", ctx, 10)
	defer ws.pubsub.Unsubscribe("aquarium:"+aquariumID.String()+":delete", ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			fmt.Fprintf(w, "event: ping\ndata: {}\n\n")
			flusher.Flush()
		case deleteFishes := <-deleteFishes:
			raw, _ := json.Marshal(deleteFishes)
			fmt.Fprintf(w, "event: fishleft\ndata: %s\n\n", raw)
			flusher.Flush()
		case fish := <-newFishes:
			raw, _ := json.Marshal(fish)
			fmt.Fprintf(w, "event: fishjoin\ndata: %s\n\n", raw)
			flusher.Flush()
		}
	}
}
