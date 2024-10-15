package webserver

import (
	"image/png"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (ws *WebServer) getFishImage(w http.ResponseWriter, r *http.Request) {
	aquariumID, err := uuid.Parse(chi.URLParam(r, "aquariumID"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 not found"))
		return
	}

	fishID, err := uuid.Parse(chi.URLParam(r, "fishID"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 not found"))
		return
	}

	// is file exist?
	img, err := ws.storage.FishImage(aquariumID, fishID)
	if err != nil {
		ws.log.Error("Failed to get fish image", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 not found"))
		return
	}

	w.Header().Set("Content-Type", "image/png")
	png.Encode(w, img)
}
