package aquarium

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	"github.com/google/uuid"
	"github.com/superbarne/fish/models"
)

type Aquarium struct {
	ID uuid.UUID

	fishes []*models.Fish

	subscribers     map[context.Context]chan *models.Fish
	subscribersLock sync.RWMutex
}

func NewAquarium(id uuid.UUID) *Aquarium {
	aquarium := &Aquarium{
		ID: id,

		subscribers: make(map[context.Context]chan *models.Fish),
	}

	// create folders
	os.MkdirAll(filepath.Join("./uploads", id.String()), os.ModePerm)
	os.MkdirAll(filepath.Join("./data", id.String()), os.ModePerm)

	// read existing fishes
	files, err := os.ReadDir(filepath.Join("./data", id.String()))
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		readFile, err := os.ReadFile(filepath.Join("./data", id.String(), file.Name()))
		if err != nil {
			panic(err)
		}

		fish := &models.Fish{}
		json.Unmarshal(readFile, fish)

		aquarium.fishes = append(aquarium.fishes, fish)
	}

	return aquarium
}

func (a *Aquarium) RealtimeFishes(ctx context.Context) <-chan *models.Fish {
	ch := make(chan *models.Fish)

	a.subscribersLock.Lock()
	defer a.subscribersLock.Unlock()
	a.subscribers[ctx] = ch

	go func() {
		<-ctx.Done()
		close(ch)
		a.subscribersLock.Lock()
		defer a.subscribersLock.Unlock()
		delete(a.subscribers, ctx)
	}()

	return ch
}

func (a *Aquarium) Fishes() []*models.Fish {
	return a.fishes
}

func (a *Aquarium) AddFish(fish *models.Fish) {
	a.fishes = append(a.fishes, fish)

	a.subscribersLock.RLock()
	defer a.subscribersLock.RUnlock()

	for _, subscriber := range a.subscribers {
		subscriber <- fish
	}
}
