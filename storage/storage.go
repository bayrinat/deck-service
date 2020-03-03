package storage

import (
	"deck-service/model"

	"github.com/google/uuid"
)

type Storage interface {
	Create(*model.Deck) (uuid.UUID, error)
	Read(uuid.UUID) *model.Deck
	Update(uuid.UUID, *model.Deck) error
	Delete(uuid.UUID) error
}
