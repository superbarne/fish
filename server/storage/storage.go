package storage

import (
	"encoding/json"
	"errors"
	"image"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/fogleman/gg"
	"github.com/google/uuid"
	"github.com/superbarne/fish/models"
)

var (
	ErrNotFound = errors.New("not found")
	ErrBadID    = errors.New("bad id")
)

// Storage is a JSON File DB
type Storage struct {
	basePath string
}

func NewStorage(basePath string) *Storage {
	return &Storage{
		basePath: basePath,
	}
}

func (s *Storage) save(path string, data interface{}) error {
	raw, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return err
	}

	if err := os.WriteFile(path, raw, os.ModePerm); err != nil {
		return err
	}

	return nil
}

// InsertAquarium inserts or updates an aquarium
func (s *Storage) InsertAquarium(aquarium *models.Aquarium) (err error) {
	if aquarium.ID == uuid.Nil {
		return ErrBadID
	}

	if aquarium.CreatedAt.IsZero() {
		aquarium.CreatedAt = time.Now()
	}
	aquarium.UpdatedAt = time.Now()

	path := filepath.Join(s.basePath, "aquariums", aquarium.ID.String(), aquarium.ID.String()+".json")
	return s.save(
		path,
		aquarium,
	)
}

// InsertFish inserts or updates a fish
func (s *Storage) InsertFish(aquariumID uuid.UUID, fish *models.Fish) (err error) {
	if aquariumID == uuid.Nil {
		return ErrBadID
	}

	if fish.ID == uuid.Nil {
		return ErrBadID
	}

	if fish.CreatedAt.IsZero() {
		fish.CreatedAt = time.Now()
	}
	fish.UpdatedAt = time.Now()

	path := filepath.Join(s.basePath, "aquariums", aquariumID.String(), "fishes", fish.ID.String()+".json")
	return s.save(
		path,
		fish,
	)
}

// Aquariums returns all aquariums
func (s *Storage) Aquariums() (aquariums []*models.Aquarium, err error) {
	aquariumsPath := filepath.Join(s.basePath, "aquariums")

	if err := os.MkdirAll(aquariumsPath, os.ModePerm); err != nil {
		return nil, err
	}

	files, err := os.ReadDir(aquariumsPath)
	if err != nil {
		return nil, err
	}

	aquariums = []*models.Aquarium{}
	for _, file := range files {
		if !file.IsDir() {
			continue
		}

		aquariumPath := filepath.Join(aquariumsPath, file.Name(), file.Name()+".json")

		raw, err := os.ReadFile(aquariumPath)
		if err != nil {
			return nil, err
		}

		metadata := &models.Aquarium{}
		if err := json.Unmarshal(raw, metadata); err != nil {
			return nil, err
		}

		aquariums = append(aquariums, metadata)
	}

	return aquariums, nil
}

// Fishes returns all fishes in an aquarium
func (s *Storage) Fishes(aquariumID uuid.UUID) (fishes []*models.Fish, err error) {
	fishesPath := filepath.Join(s.basePath, "aquariums", aquariumID.String(), "fishes")

	files, err := os.ReadDir(fishesPath)
	if err != nil {
		return nil, err
	}

	fishes = []*models.Fish{}
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		raw, err := os.ReadFile(filepath.Join(fishesPath, file.Name()))
		if err != nil {
			return nil, err
		}

		fish := &models.Fish{}
		if err := json.Unmarshal(raw, fish); err != nil {
			return nil, err
		}

		fishes = append(fishes, fish)
	}

	return fishes, nil
}

// Aquarium returns an aquarium
func (s *Storage) Aquarium(aquariumID uuid.UUID) (aquarium *models.Aquarium, err error) {
	if aquariumID == uuid.Nil {
		return nil, ErrBadID
	}

	aquariumPath := filepath.Join(s.basePath, "aquariums", aquariumID.String(), aquariumID.String()+".json")

	raw, err := os.ReadFile(aquariumPath)
	if err != nil {
		return nil, err
	}

	aquarium = &models.Aquarium{}
	if err := json.Unmarshal(raw, aquarium); err != nil {
		return nil, err
	}

	return aquarium, nil
}

// Fish returns a fish
func (s *Storage) Fish(aquariumID uuid.UUID, fishID uuid.UUID) (fish *models.Fish, err error) {
	if aquariumID == uuid.Nil {
		return nil, ErrBadID
	}

	if fishID == uuid.Nil {
		return nil, ErrBadID
	}

	fishPath := filepath.Join(s.basePath, "aquariums", aquariumID.String(), "fishes", fishID.String()+".json")

	raw, err := os.ReadFile(fishPath)
	if err != nil {
		return nil, err
	}

	fish = &models.Fish{}
	if err := json.Unmarshal(raw, fish); err != nil {
		return nil, err
	}

	return fish, nil
}

// FishImagePath returns a fish image path
func (s *Storage) FishImagePath(aquariumID uuid.UUID, fishID uuid.UUID) (path string, err error) {
	if aquariumID == uuid.Nil {
		return "", ErrBadID
	}

	if fishID == uuid.Nil {
		return "", ErrBadID
	}

	fishPath := filepath.Join(s.basePath, "aquariums", aquariumID.String(), "fishes_images")
	filename := fishID.String() + ".png"

	return filepath.Join(fishPath, filename), nil
}

// FishImage returns a fish image
func (s *Storage) FishImage(aquariumID uuid.UUID, fishID uuid.UUID) (img image.Image, err error) {
	path, err := s.FishImagePath(aquariumID, fishID)
	if err != nil {
		return nil, err
	}

	img, err = gg.LoadImage(path)
	if err != nil {
		return nil, err
	}

	return img, nil
}

// DeleteAquarium deletes an aquarium
func (s *Storage) DeleteAquarium(id uuid.UUID) (err error) {
	if id == uuid.Nil {
		return ErrBadID
	}

	aquariumPath := filepath.Join(s.basePath, "aquariums", id.String())

	if err := os.RemoveAll(aquariumPath); err != nil {
		return err
	}

	return nil
}

// DeleteFish deletes a fish
func (s *Storage) DeleteFish(aquariumID uuid.UUID, fishID uuid.UUID) (err error) {
	if aquariumID == uuid.Nil {
		return ErrBadID
	}

	if fishID == uuid.Nil {
		return ErrBadID
	}

	fishPath := filepath.Join(s.basePath, "aquariums", aquariumID.String(), "fishes", fishID.String()+".json")

	if err := os.Remove(fishPath); err != nil {
		return err
	}

	// delete image
	fishImagePath, err := s.FishImagePath(aquariumID, fishID)
	if err != nil {
		return err
	}

	if err := os.Remove(fishImagePath); err != nil {
		return err
	}

	return nil
}

func (s *Storage) SaveTmpFishImageFromRequest(aquariumID uuid.UUID, fishID uuid.UUID, file multipart.File, multipartHeader *multipart.FileHeader) (string, error) {
	if aquariumID == uuid.Nil {
		return "", ErrBadID
	}

	if fishID == uuid.Nil {
		return "", ErrBadID
	}

	tmpFolder := filepath.Join(os.TempDir(), "aquariums", aquariumID.String())
	if err := os.MkdirAll(tmpFolder, os.ModePerm); err != nil {
		return "", err
	}

	fileName := fishID.String() + filepath.Ext(multipartHeader.Filename)

	out, err := os.Create(filepath.Join(tmpFolder, fileName))
	if err != nil {
		return "", err
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		return "", err
	}

	return filepath.Join(tmpFolder, fileName), nil
}
