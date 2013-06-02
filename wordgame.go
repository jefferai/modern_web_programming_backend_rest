package main

import (
	"crypto/rand"
	"github.com/coocood/jas"
	"math"
	"math/big"
	"strings"
)

// STARTRESULTSTRUCT OMIT
type Result struct {
	Id            int64  `json:"id"`
	NumGuesses    int64  `json:"numguesses"`
	CurrentString string `json:"currentstring"`
	Correct       bool   `json:"correct"`
}
// ENDRESULTSTRUCT OMIT

// STARTNEWVARS OMIT
type Game struct {
	NumGuesses int64
	Word       string
	MaskedWord string
	Name       string
}

var (
	games          = map[int64]Game{}
	words          = []string{"aardvark", "labradoodle", "kittycat", "porpoise", "brontosaurus"}
	underscoreFunc = func(r rune) rune { return '_' }
)
// ENDNEWVARS OMIT

// STARTWGDEF OMIT
type Wordgame struct{}

func (*Wordgame) GetNewgame(ctx *jas.Context) {
	// ENDWGDEF OMIT

	// STARTIEDEF OMIT
	randint, err := rand.Int(rand.Reader, big.NewInt(int64(math.MaxUint32)))

	if err != nil {
		ctx.Error = jas.NewInternalError("Got an error generating a game ID")
		panic(ctx.Error)
	}
	// ENDIEDEF OMIT

	// STARTREDEF OMIT
	// Commented out for example reasons; will panic with default error if not found
	// name := ctx.RequireString("name")

	name, err := ctx.FindString("name")
	if err != nil {
		ctx.Error = jas.NewRequestError("Player's name must be provided")
		panic(ctx.Error)
	}
	// ENDREDEF OMIT

	currgame := randint.Int64()
	wordnum := currgame % int64(len(words))

	games[currgame] = Game{
		NumGuesses: 0,
		Word:       words[wordnum],
		MaskedWord: strings.Map(underscoreFunc, words[wordnum]),
		Name:       name,
	}

	result := Result{
		Id:            currgame,
		NumGuesses:    0,
		CurrentString: games[currgame].MaskedWord,
	}

	ctx.Data = &result
	return
}
