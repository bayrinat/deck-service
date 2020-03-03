package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"deck-service/model"
	"deck-service/storage"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type handlers struct {
	storage storage.Storage
	logger  *zap.Logger
}

func New(storage storage.Storage, logger *zap.Logger) *handlers {
	return &handlers{storage, logger}
}

func (h *handlers) CreateDeckHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Request received to create a new deck")

	var (
		shuffled bool
	)
	shuffledVal, ok := r.URL.Query()["shuffled"]
	if ok && strings.EqualFold(shuffledVal[0], "true") {
		shuffled = true
	}
	cards, _ := r.URL.Query()["cards"]
	deck, err := model.NewDeck(cards, shuffled)
	if err != nil {
		h.logger.Error("Failed to create new deck", zap.Error(err))
		h.send(w, http.StatusBadRequest, nil)
		return
	}
	id, err := h.storage.Create(deck)
	if err != nil {
		h.logger.Error("Failed to save a deck in storage", zap.Error(err))
		h.send(w, http.StatusInternalServerError, nil)
		return
	}

	h.logger.Info("Successfully created the deck ", zap.String("id", id.String()))
	h.send(w, http.StatusCreated, &struct {
		ID        string `json:"deck_id"`
		Shuffled  bool   `json:"shuffled"`
		Remaining int    `json:"remaining"`
	}{
		id.String(),
		deck.Shuffled,
		deck.Remaining,
	})
}

func (h *handlers) OpenDeckHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	h.logger.Info("Request received to open the deck by id: ", zap.String("id", id))

	uuid, err := uuid.Parse(id)
	if err != nil {
		h.logger.Error("Failed to transform ID to UUID", zap.Error(err), zap.String("id", id))
		h.send(w, http.StatusBadRequest, nil)
		return
	}

	deck := h.storage.Read(uuid)
	if deck == nil {
		h.logger.Error("Failed to get deck by ID", zap.String("id", id))
		h.send(w, http.StatusNotFound, nil)
		return
	}
	h.logger.Info("Successfully got the deck ", zap.String("id", id))
	h.send(w, http.StatusOK, &struct {
		ID        string        `json:"deck_id"`
		Shuffled  bool          `json:"shuffled"`
		Remaining int           `json:"remaining"`
		Cards     []*model.Card `json:"cards"`
	}{
		id,
		deck.Shuffled,
		deck.Remaining,
		deck.Cards,
	})
	return
}

func (h *handlers) DrawCard(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	h.logger.Info("Request received to draw cards of the deck with id: ", zap.String("id", id))

	uuid, err := uuid.Parse(id)
	if err != nil {
		h.logger.Error("Failed to transform ID to UUID", zap.Error(err), zap.String("id", id))
		h.send(w, http.StatusBadRequest, nil)
		return
	}

	deck := h.storage.Read(uuid)
	if deck == nil {
		h.logger.Error("Failed to get deck by ID", zap.String("id", id))
		h.send(w, http.StatusNotFound, nil)
		return
	}
	countVal, ok := r.URL.Query()["count"]
	if !ok {
		h.logger.Error("query parameter `count` must be defined")
		h.send(w, http.StatusBadRequest, nil)
		return
	}
	var count int
	count, err = strconv.Atoi(countVal[0])
	if err != nil {
		h.logger.Error("Failed to convert `count` to number", zap.Error(err))
		h.send(w, http.StatusBadRequest, nil)
		return
	}
	if count <= 0 {
		h.logger.Error("Failed to draw cards, count parameter must be positive", zap.Int("count", count))
		h.send(w, http.StatusBadRequest, nil)
		return
	}
	cards, err := deck.Draw(count)
	if err != nil {
		h.logger.Error("Failed to draw cards", zap.Error(err))
		h.send(w, http.StatusBadRequest, nil)
		return
	}

	err = h.storage.Update(uuid, deck)
	if err != nil {
		h.logger.Error("Failed to convert `count` to number", zap.Error(err))
		h.send(w, http.StatusBadRequest, nil)
		return
	}

	h.logger.Info("Successfully got the deck ", zap.String("id", id))
	h.send(w, http.StatusOK, &struct {
		Cards []*model.Card `json:"cards"`
	}{
		cards,
	})
	return
}

func (h *handlers) send(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	if data == nil {
		return
	}
	w.Header().Add("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		h.logger.Error("Failed to encode data", zap.Error(err))
	}
}
