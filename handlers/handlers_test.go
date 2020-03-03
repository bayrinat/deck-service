package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"deck-service/model"
	"deck-service/storage"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type knownDeck struct {
	id   string
	deck *model.Deck
}

var knownDecks []*knownDeck

func testHandlers() *handlers {
	storage := storage.NewConcurrentMemory()
	logger := zap.NewNop()

	deck, _ := model.NewDeck(nil, true)
	id, _ := storage.Create(deck)
	knownDecks = append(knownDecks, &knownDeck{id.String(), deck})

	deck, _ = model.NewDeck(nil, false)
	id, _ = storage.Create(deck)
	knownDecks = append(knownDecks, &knownDeck{id.String(), deck})

	return &handlers{storage, logger}
}

func TestHandlers(t *testing.T) {
	handlers := testHandlers()

	tests := []struct {
		name       string
		method     string
		url        string
		path       string
		handler    http.HandlerFunc
		checkBody  func(*bytes.Buffer) bool
		wantStatus int
	}{
		// TODO(bayrinat): add negative scenarios
		{
			name:       "Create deck",
			method:     http.MethodPost,
			url:        "/deck?shuffled=true",
			path:       "/deck",
			handler:    handlers.CreateDeckHandler,
			wantStatus: http.StatusCreated,
			checkBody: func(body *bytes.Buffer) bool {
				expected := struct {
					ID        string `json:"deck_id"`
					Shuffled  bool   `json:"shuffled"`
					Remaining int    `json:"remaining"`
				}{}
				err := json.Unmarshal(body.Bytes(), &expected)
				if err != nil {
					return false
				}
				_, err = uuid.Parse(expected.ID)
				if err != nil {
					return false
				}
				if !expected.Shuffled {
					return false
				}
				if expected.Remaining != 52 {
					return false
				}
				return true
			},
		},
		{
			name:       "Open deck",
			method:     http.MethodGet,
			url:        "/decks/" + knownDecks[1].id,
			path:       "/decks/{id}",
			handler:    handlers.OpenDeckHandler,
			wantStatus: http.StatusOK,
			checkBody: func(body *bytes.Buffer) bool {
				expected := struct {
					ID        string        `json:"deck_id"`
					Shuffled  bool          `json:"shuffled"`
					Remaining int           `json:"remaining"`
					Cards     []*model.Card `json:"cards"`
				}{}
				err := json.Unmarshal(body.Bytes(), &expected)
				if err != nil {
					return false
				}
				_, err = uuid.Parse(expected.ID)
				if err != nil {
					return false
				}
				if expected.Shuffled {
					return false
				}
				if expected.Remaining != 52 || len(expected.Cards) != 52 {
					return false
				}
				if !reflect.DeepEqual(expected.Cards, model.DefaultDeckCards) {
					return false
				}
				return true
			},
		},
		{
			name:       "Try to open nonexistent deck",
			method:     http.MethodGet,
			url:        "/decks/" + uuid.New().String(),
			path:       "/decks/{id}",
			handler:    handlers.OpenDeckHandler,
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "Draw cards",
			method:     http.MethodPost,
			url:        "/decks/" + knownDecks[1].id + "/draw?count=10",
			path:       "/decks/{id}/draw",
			handler:    handlers.DrawCard,
			wantStatus: http.StatusOK,
			checkBody: func(body *bytes.Buffer) bool {
				expected := struct {
					Cards []*model.Card `json:"cards"`
				}{}
				err := json.Unmarshal(body.Bytes(), &expected)
				if err != nil {
					return false
				}
				if len(expected.Cards) != 10 {
					return false
				}
				if !reflect.DeepEqual(expected.Cards, model.DefaultDeckCards[:10]) {
					return false
				}
				return true
			},
		},
		{
			name:       "Draw cards on nonexistent deck",
			method:     http.MethodPost,
			url:        "/decks/" + uuid.New().String() + "/draw",
			path:       "/decks/{id}/draw",
			handler:    handlers.DrawCard,
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "Draw cards without count",
			method:     http.MethodPost,
			url:        "/decks/" + knownDecks[0].id + "/draw",
			path:       "/decks/{id}/draw",
			handler:    handlers.DrawCard,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, tt.url, nil)
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			router := mux.NewRouter()
			router.HandleFunc(tt.path, tt.handler)
			router.ServeHTTP(rr, req)

			if rr.Code != tt.wantStatus {
				t.Errorf("Wrong status got = %v, want %v", rr.Code, tt.wantStatus)
			}
			if tt.checkBody != nil && !tt.checkBody(rr.Body) {
				t.Errorf("Not expected answer")
			}
		})
	}
}
