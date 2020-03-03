package model

import "strings"

type Value string

var Values = map[string]Value{
	"two":   "2",
	"three": "3",
	"four":  "4",
	"five":  "5",
	"six":   "6",
	"seven": "7",
	"eight": "8",
	"nine":  "9",
	"ten":   "10",
	"jack":  "j",
	"queen": "q",
	"king":  "k",
	"ace":   "a",
}

func validateValue(name string) (Value, bool) {
	v, ok := Values[strings.ToLower(name)]
	return v, ok
}
