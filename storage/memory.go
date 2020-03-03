package storage

import (
	"deck-service/model"

	"github.com/google/uuid"
)

type memory struct {
	decks map[uuid.UUID]*model.Deck
}

func (m *memory) Create(deck *model.Deck) (uuid.UUID, error) {
	id := uuid.New()
	m.decks[id] = deck
	return id, nil
}

func (m *memory) Read(id uuid.UUID) *model.Deck {
	return m.decks[id]
}

func (m *memory) Update(id uuid.UUID, deck *model.Deck) error {
	m.decks[id] = deck
	return nil
}

func (m *memory) Delete(id uuid.UUID) error {
	delete(m.decks, id)
	return nil
}
