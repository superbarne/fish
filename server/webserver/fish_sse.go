package webserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/superbarne/fish/aquarium"
)

func (ws *WebServer) sseFish(w http.ResponseWriter, r *http.Request) {
	// validate id
	aquariumID, err := uuid.Parse(chi.URLParam(r, "aquariumID"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 not found"))
		return
	}

	var aquarium *aquarium.Aquarium
	func() {
		ws.aquariumsLock.RLock()
		defer ws.aquariumsLock.RUnlock()

		var ok bool
		if aquarium, ok = ws.aquariums[aquariumID]; !ok {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 not found"))
			return
		}
	}()

	if aquarium == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 not found"))
		return
	}

	ctx := r.Context()
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	f, ok := w.(http.Flusher)
	if !ok {
		return
	}

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	fmt.Fprintf(w, "event: ping\ndata: {}\n\n")

	// send old fishes
	for _, fish := range aquarium.Fishes() {
		raw, _ := json.Marshal(fish)
		fmt.Fprintf(w, "event: fish\ndata: %s\n\n", raw)
	}
	f.Flush()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			fmt.Fprintf(w, "event: ping\ndata: {}\n\n")
			f.Flush()
		case fish := <-aquarium.RealtimeFishes(ctx):
			raw, _ := json.Marshal(fish)
			fmt.Fprintf(w, "event: fish\ndata: %s\n\n", raw)
			f.Flush()
		}
	}
}
