package model

import (
	"fmt"
	"strings"
)

type Card struct {
	// TODO(bayrinat): validate value and suite
	Value Value  `json:"value"`
	Suit  Suit   `json:"suit"`
	Code  string `json:"code"`
}

func NewCard(value Value, suit Suit) (c *Card) {
	// TODO(bayrinat): it is easy to mix up value and suit, fix it!
	// TODO(bayrinat): validate value and suite, return error instead of nil
	v, okVal := validateValue(string(value))
	s, okSuit := validateSuit(string(suit))
	if okVal && okSuit {
		c = &Card{
			Value: v,
			Suit:  s,
		}
		c.Code = c.makeCode()
	}
	return
}

func (c *Card) makeCode() string {
	return fmt.Sprintf("%s%s", c.Value, c.Suit)
}

func GetCardByCode(code string) *Card {
	// NOTE(bayrinat): linear search is good enough here
	for _, card := range validCards {
		if strings.EqualFold(code, card.makeCode()) {
			return card
		}
	}
	return nil
}

var validCards = []*Card{
	NewCard("Two", "Clubs"),
	NewCard("Three", "Clubs"),
	NewCard("Four", "Clubs"),
	NewCard("Five", "Clubs"),
	NewCard("Six", "Clubs"),
	NewCard("Seven", "Clubs"),
	NewCard("Eight", "Clubs"),
	NewCard("Nine", "Clubs"),
	NewCard("Ten", "Clubs"),
	NewCard("Jack", "Clubs"),
	NewCard("Queen", "Clubs"),
	NewCard("King", "Clubs"),
	NewCard("Ace", "Clubs"),

	NewCard("Two", "Diamonds"),
	NewCard("Three", "Diamonds"),
	NewCard("Four", "Diamonds"),
	NewCard("Five", "Diamonds"),
	NewCard("Six", "Diamonds"),
	NewCard("Seven", "Diamonds"),
	NewCard("Eight", "Diamonds"),
	NewCard("Nine", "Diamonds"),
	NewCard("Ten", "Diamonds"),
	NewCard("Jack", "Diamonds"),
	NewCard("Queen", "Diamonds"),
	NewCard("King", "Diamonds"),
	NewCard("Ace", "Diamonds"),

	NewCard("Two", "Hearts"),
	NewCard("Three", "Hearts"),
	NewCard("Four", "Hearts"),
	NewCard("Five", "Hearts"),
	NewCard("Six", "Hearts"),
	NewCard("Seven", "Hearts"),
	NewCard("Eight", "Hearts"),
	NewCard("Nine", "Hearts"),
	NewCard("Ten", "Hearts"),
	NewCard("Jack", "Hearts"),
	NewCard("Queen", "Hearts"),
	NewCard("King", "Hearts"),
	NewCard("Ace", "Hearts"),

	NewCard("Two", "Spades"),
	NewCard("Three", "Spades"),
	NewCard("Four", "Spades"),
	NewCard("Five", "Spades"),
	NewCard("Six", "Spades"),
	NewCard("Seven", "Spades"),
	NewCard("Eight", "Spades"),
	NewCard("Nine", "Spades"),
	NewCard("Ten", "Spades"),
	NewCard("Jack", "Spades"),
	NewCard("Queen", "Spades"),
	NewCard("King", "Spades"),
	NewCard("Ace", "Spades"),
}
