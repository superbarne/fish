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
	ErrBadFishID     = errors.New("bad fish ID")
	ErrBadAquariumID = errors.New("bad aquarium ID")
	ErrNotFound      = errors.New("not found")
)

type Storage struct {
	basePath string
}

func NewStorage(basePath string) *Storage {
	return &Storage{
		basePath: basePath,
	}
}

func (s *Storage) SaveAquarium(metadata *models.Aquarium) error {
	if metadata.ID == uuid.Nil {
		return ErrBadAquariumID
	}

	raw, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	aquariumsPath := filepath.Join(s.basePath, "aquariums", metadata.ID.String())
	filename := metadata.ID.String() + ".json"

	metadata.UpdatedAt = time.Now()
	// exists?
	if _, err := os.Stat(aquariumsPath); os.IsNotExist(err) {
		metadata.CreatedAt = time.Now()
	}

	// create folder
	if err := os.MkdirAll(aquariumsPath, os.ModePerm); err != nil {
		return err
	}

	// save to file
	if err := os.WriteFile(filepath.Join(aquariumsPath, filename), raw, os.ModePerm); err != nil {
		return err
	}

	return nil
}

func (s *Storage) Aquarium(id uuid.UUID) (*models.Aquarium, error) {
	if id == uuid.Nil {
		return nil, ErrBadAquariumID
	}

	aquariumsPath := filepath.Join(s.basePath, "aquariums", id.String())
	filename := id.String() + ".json"

	raw, err := os.ReadFile(filepath.Join(aquariumsPath, filename))
	if err != nil {
		return nil, err
	}

	metadata := &models.Aquarium{}
	if err := json.Unmarshal(raw, metadata); err != nil {
		return nil, err
	}

	return metadata, nil
}

func (s *Storage) DeleteAquarium(id uuid.UUID) error {
	if id == uuid.Nil {
		return ErrBadAquariumID
	}

	aquariumsPath := filepath.Join(s.basePath, "aquariums", id.String())

	if err := os.RemoveAll(aquariumsPath); err != nil {
		return err
	}

	return nil
}

func (s *Storage) SaveFishMetadata(aquariumID uuid.UUID, fish *models.Fish) error {
	if fish.ID == uuid.Nil {
		return ErrBadFishID
	}

	if aquariumID == uuid.Nil || fish.AquariumID != aquariumID {
		return ErrBadAquariumID
	}

	raw, err := json.Marshal(fish)
	if err != nil {
		return err
	}

	fishesPath := filepath.Join(s.basePath, "aquariums", aquariumID.String(), "fishes")
	filename := fish.ID.String() + ".json"
	fish.UpdatedAt = time.Now()

	// exists?
	if _, err := os.Stat(fishesPath); os.IsNotExist(err) {
		fish.CreatedAt = time.Now()
	}

	// create folder
	if err := os.MkdirAll(fishesPath, os.ModePerm); err != nil {
		return err
	}

	// save to file
	if err := os.WriteFile(filepath.Join(fishesPath, filename), raw, os.ModePerm); err != nil {
		return err
	}

	return nil
}

func (s *Storage) FishMetadata(aquariumID, fishID uuid.UUID) (*models.Fish, error) {
	if fishID == uuid.Nil {
		return nil, ErrBadFishID
	}

	if aquariumID == uuid.Nil {
		return nil, ErrBadAquariumID
	}

	fishesPath := filepath.Join(s.basePath, "aquariums", aquariumID.String(), "fishes")
	filename := fishID.String() + ".json"

	raw, err := os.ReadFile(filepath.Join(fishesPath, filename))
	if err != nil {
		return nil, err
	}

	fish := &models.Fish{}
	if err := json.Unmarshal(raw, fish); err != nil {
		return nil, err
	}

	return fish, nil
}

func (s *Storage) SaveTmpFishImageFromRequest(aquariumID uuid.UUID, fishID uuid.UUID, file multipart.File, multipartHeader *multipart.FileHeader) (string, error) {
	if aquariumID == uuid.Nil {
		return "", ErrBadAquariumID
	}

	if fishID == uuid.Nil {
		return "", ErrBadFishID
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

func (s *Storage) FishImagePath(aquariumID, fishID uuid.UUID) string {
	fishesPath := filepath.Join(s.basePath, "aquariums", aquariumID.String(), "fishes_images")
	filename := fishID.String() + ".png"

	return filepath.Join(fishesPath, filename)
}

func (s *Storage) FishImage(aquariumID, fishID uuid.UUID) (image.Image, error) {
	if fishID == uuid.Nil {
		return nil, ErrBadFishID
	}

	if aquariumID == uuid.Nil {
		return nil, ErrBadAquariumID
	}

	fishesPath := filepath.Join(s.basePath, "aquariums", aquariumID.String(), "fishes_images")
	filename := fishID.String() + ".png"

	img, err := gg.LoadImage(filepath.Join(fishesPath, filename))
	if err != nil {
		return nil, err
	}

	return img, nil
}

func (s *Storage) DeleteFish(aquariumID, fishID uuid.UUID) error {
	if fishID == uuid.Nil {
		return ErrBadFishID
	}

	if aquariumID == uuid.Nil {
		return ErrBadAquariumID
	}

	fishesPath := filepath.Join(s.basePath, "aquariums", aquariumID.String(), "fishes")
	filename := fishID.String() + ".json"

	if err := os.Remove(filepath.Join(fishesPath, filename)); err != nil {
		return err
	}

	fishesPath = filepath.Join(s.basePath, "aquariums", aquariumID.String(), "fishes_images")
	filename = fishID.String() + ".png"

	if err := os.Remove(filepath.Join(fishesPath, filename)); err != nil {
		return err
	}

	return nil
}

func (s *Storage) Fishes(aquariumID uuid.UUID) ([]*models.Fish, error) {
	if aquariumID == uuid.Nil {
		return nil, ErrBadAquariumID
	}

	fishesPath := filepath.Join(s.basePath, "aquariums", aquariumID.String(), "fishes")

	if err := os.MkdirAll(fishesPath, os.ModePerm); err != nil {
		return nil, err
	}

	files, err := os.ReadDir(fishesPath)
	if err != nil {
		return nil, err
	}

	fishes := []*models.Fish{}
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
