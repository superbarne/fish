package aquarium

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/superbarne/fish/models"
	"github.com/superbarne/fish/storage"
)

type Aquarium struct {
	ID uuid.UUID

	fishes []*models.Fish

	storage         *storage.Storage
	subscribers     map[context.Context]chan *models.Fish
	subscribersLock sync.RWMutex
}

func NewAquarium(id uuid.UUID, storage *storage.Storage) *Aquarium {
	aquarium := &Aquarium{
		ID: id,

		storage:     storage,
		subscribers: make(map[context.Context]chan *models.Fish),
	}

	fishes, err := aquarium.storage.Fishes(id)
	if err != nil {
		panic(err)
	}

	aquarium.fishes = append(aquarium.fishes, fishes...)

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
