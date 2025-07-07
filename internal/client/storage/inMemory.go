package storage

import (
	"errors"
	"fmt"
	"github.com/Fuonder/datakeeper.git/internal/models"
	"sync"
)

type InMemoryStorage struct {
	AuthToken string
	Data      models.ObjectList
	mu        sync.RWMutex
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		AuthToken: "",
		Data:      models.ObjectList{},
	}
}

func (s *InMemoryStorage) SetToken(token string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.AuthToken = token
}

func (s *InMemoryStorage) GetToken() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.AuthToken
}

// AddItem adds or replaces an item in the appropriate list.
func (s *InMemoryStorage) AddItem(item interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	switch v := item.(type) {
	case models.LoginData:
		for i, existing := range s.Data.LoginObjects {
			if existing.ID == v.ID {
				s.Data.LoginObjects[i] = v
				return nil
			}
		}
		s.Data.LoginObjects = append(s.Data.LoginObjects, v)
	case models.TextData:
		for i, existing := range s.Data.TextObjects {
			if existing.ID == v.ID {
				s.Data.TextObjects[i] = v
				return nil
			}
		}
		s.Data.TextObjects = append(s.Data.TextObjects, v)
	case models.CreditCardData:
		for i, existing := range s.Data.CreditCardObjects {
			if existing.ID == v.ID {
				s.Data.CreditCardObjects[i] = v
				return nil
			}
		}
		s.Data.CreditCardObjects = append(s.Data.CreditCardObjects, v)
	case models.FileData:
		for i, existing := range s.Data.FileObjects {
			if existing.ID == v.ID {
				s.Data.FileObjects[i] = v
				return nil
			}
		}
		s.Data.FileObjects = append(s.Data.FileObjects, v)
	default:
		return fmt.Errorf("unsupported item type: %T", item)
	}

	return nil
}

// Get methods
func (s *InMemoryStorage) GetLoginItemByIndex(idx int) (models.LoginData, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if idx < 0 || idx >= len(s.Data.LoginObjects) {
		return models.LoginData{}, errors.New("index out of range")
	}
	return s.Data.LoginObjects[idx], nil
}

func (s *InMemoryStorage) GetTextItemByIndex(idx int) (models.TextData, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if idx < 0 || idx >= len(s.Data.TextObjects) {
		return models.TextData{}, errors.New("index out of range")
	}
	return s.Data.TextObjects[idx], nil
}

func (s *InMemoryStorage) GetCardItemByIndex(idx int) (models.CreditCardData, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if idx < 0 || idx >= len(s.Data.CreditCardObjects) {
		return models.CreditCardData{}, errors.New("index out of range")
	}
	return s.Data.CreditCardObjects[idx], nil
}

func (s *InMemoryStorage) GetFileItemByIndex(idx int) (models.FileData, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if idx < 0 || idx >= len(s.Data.FileObjects) {
		return models.FileData{}, errors.New("index out of range")
	}
	return s.Data.FileObjects[idx], nil
}
