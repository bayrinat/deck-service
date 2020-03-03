package storage

import (
	"sync"

	"deck-service/model"

	"github.com/google/uuid"
)

type concurrentMemory struct {
	mu    *sync.Mutex
	decks map[uuid.UUID]*model.Deck
}

func NewConcurrentMemory() Storage {
	return &concurrentMemory{
		mu:    &sync.Mutex{},
		decks: make(map[uuid.UUID]*model.Deck),
	}
}

func (m *concurrentMemory) Create(deck *model.Deck) (uuid.UUID, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	id := uuid.New()
	m.decks[id] = deck
	return id, nil
}

func (m *concurrentMemory) Read(id uuid.UUID) *model.Deck {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.decks[id]
}

func (m *concurrentMemory) Update(id uuid.UUID, deck *model.Deck) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.decks[id] = deck
	return nil
}

func (m *concurrentMemory) Delete(id uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.decks, id)
	return nil
}
