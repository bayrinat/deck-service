package model

import "strings"

type Suit string

var Suits = map[string]Suit{
	"clubs":    "c",
	"diamonds": "d",
	"hearts":   "h",
	"spades":   "s",
}

func validateSuit(name string) (Suit, bool) {
	s, ok := Suits[strings.ToLower(name)]
	return s, ok
}
