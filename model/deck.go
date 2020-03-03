package model

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"time"
)

var DefaultDeckCards = validCards

var ErrInvalidCount = errors.New("invalid count")
var ErrAlreadyDrawn = errors.New("the whole deck was already drawn")

type Deck struct {
	Shuffled  bool    `json:"shuffled"`
	Remaining int     `json:"remaining"`
	Cards     []*Card `json:"cards"`
}

func NewDeck(cards []string, shuffled bool) (*Deck, error) {
	rand.Seed(time.Now().Unix())
	deck := Deck{
		Shuffled: shuffled,
	}
	// collect deck
	if len(cards) == 0 {
		deck.Cards = make([]*Card, len(DefaultDeckCards))
		copy(deck.Cards, DefaultDeckCards)
	}
	if len(cards) > 0 {
		deck.Cards = make([]*Card, len(cards))
		for i, cardCode := range cards {
			card := GetCardByCode(cardCode)
			if card == nil {
				return nil, fmt.Errorf("invalid card's makeCode %s", cardCode)
			}
			deck.Cards[i] = card
		}
	}
	deck.Remaining = len(deck.Cards)
	// shuffle deck
	if shuffled {
		deck.shuffle()
	}
	return &deck, nil
}

func (d *Deck) Draw(count int) ([]*Card, error) {
	if d.Remaining == 0 {
		return nil, ErrAlreadyDrawn
	}
	if count <= 0 {
		return nil, ErrInvalidCount
	}
	// draw all Cards requested or all Cards till the end of the deck
	remaining := int(math.Max(0, float64(d.Remaining-count)))
	newIndex := len(d.Cards) - remaining
	oldIndex := len(d.Cards) - d.Remaining
	d.Remaining = remaining
	return d.Cards[oldIndex:newIndex], nil
}

func (d *Deck) shuffle() {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for n := len(d.Cards); n > 0; n-- {
		randIndex := r.Intn(n)
		d.Cards[n-1], d.Cards[randIndex] = d.Cards[randIndex], d.Cards[n-1]
	}
}
