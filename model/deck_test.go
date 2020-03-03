package model

import (
	"reflect"
	"testing"
)

func TestNewDeck(t *testing.T) {
	type args struct {
		cards    []string
		shuffled bool
	}
	tests := []struct {
		name    string
		args    args
		want    *Deck
		wantErr bool
	}{
		{
			name: "Normal with given cards",
			args: args{
				cards:    []string{"AS", "KH", "8C"},
				shuffled: false,
			},
			want: &Deck{
				Shuffled:  false,
				Remaining: 3,
				Cards: []*Card{
					NewCard("Ace", "Spades"),
					NewCard("King", "Hearts"),
					NewCard("Eight", "Clubs"),
				},
			},
			wantErr: false,
		},
		{
			name: "Bad makeCode with given cards",
			args: args{
				cards:    []string{"AS", "KH", "AA"},
				shuffled: false,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Default cards",
			args: args{
				cards:    nil,
				shuffled: false,
			},
			want: &Deck{
				Shuffled:  false,
				Remaining: 52,
				Cards:     DefaultDeckCards,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewDeck(tt.args.cards, tt.args.shuffled)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDeck() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDeck() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_deckImpl_Draw(t *testing.T) {
	type fields struct {
		Shuffled  bool
		Remaining int
		Cards     []*Card
	}
	type args struct {
		count int
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		want          []*Card
		wantRemaining int
		wantErr       bool
	}{
		{
			name: "Draw some",
			fields: fields{
				Shuffled:  false,
				Remaining: 52,
				Cards:     DefaultDeckCards,
			},
			args: args{
				count: 3,
			},
			want: []*Card{
				NewCard("Two", "Clubs"),
				NewCard("Three", "Clubs"),
				NewCard("Four", "Clubs"),
			},
			wantRemaining: 49,
			wantErr:       false,
		},
		{
			name: "Draw more than capacity",
			fields: fields{
				Shuffled:  false,
				Remaining: 52,
				Cards:     DefaultDeckCards,
			},
			args: args{
				count: 93,
			},
			want:          DefaultDeckCards,
			wantRemaining: 0,
			wantErr:       false,
		},
		{
			name: "Draw when a deck was already drawn",
			fields: fields{
				Shuffled:  false,
				Remaining: 0,
				Cards:     DefaultDeckCards,
			},
			args: args{
				count: 10,
			},
			want:          nil,
			wantRemaining: 0,
			wantErr:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Deck{
				Shuffled:  tt.fields.Shuffled,
				Remaining: tt.fields.Remaining,
				Cards:     tt.fields.Cards,
			}
			got, err := d.Draw(tt.args.count)
			if (err != nil) != tt.wantErr {
				t.Errorf("Draw() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Draw() got = %v, want %v", got, tt.want)
			}
			if d.Remaining != tt.wantRemaining {
				t.Errorf("Draw() remaining got = %v, want %v", d.Remaining, tt.wantRemaining)
			}
		})
	}
}

func Test_deckImpl_shuffle(t *testing.T) {
	type fields struct {
		Shuffled  bool
		Remaining int
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "Shuffle!!!",
			fields: fields{
				Shuffled:  true,
				Remaining: 52,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Deck{
				Shuffled:  tt.fields.Shuffled,
				Remaining: tt.fields.Remaining,
			}
			copy(d.Cards, DefaultDeckCards)
			d.shuffle()
			// TODO(bayrinat): unshuffled deck is valid even after shuffle, fix this check
			if reflect.DeepEqual(d.Cards, DefaultDeckCards) {
				t.Errorf("Shuffle() must shuffle")
			}
			if d.Remaining != tt.fields.Remaining {
				t.Errorf("Shuffle() remaining got = %v, want %v", d.Remaining, tt.fields.Remaining)
			}
		})
	}
}
