package main

import (
	"crypto/rand"
	"github.com/coocood/jas"
	"math"
	"math/big"
	"strings"
)

// This ignores safety for simplicity
func UnhideByte(guess byte, word string, maskedword string) string {
	wordbytes := []byte(word)
	maskedwordbytes := []byte(maskedword)
	for pos := range wordbytes {
		if wordbytes[pos] == guess {
			maskedwordbytes[pos] = wordbytes[pos]
		}
	}
	return string(maskedwordbytes)
}

type Result struct {
	Id            int64  `json:"id"`
	NumGuesses    int64  `json:"numguesses"`
	CurrentString string `json:"currentstring"`
	Correct       bool   `json:"correct"`
}

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

type Wordgame struct{}

func (*Wordgame) GetNewgame(ctx *jas.Context) {

	randint, err := rand.Int(rand.Reader, big.NewInt(int64(math.MaxUint32)))

	if err != nil {
		ctx.Error = jas.NewInternalError("Got an error generating a game ID")
		panic(nil)
	}

	name, err := ctx.FindString("name")
	if err != nil {
		ctx.Error = jas.NewRequestError("Player's name must be provided")
		panic(nil)
	}

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

func (*Wordgame) PostGuess(ctx *jas.Context) {
	/*
	As a reminder, these are the relevant data structures:

	type ReceivedData struct {
		Id       int64  `json:"id"`
		Guess    string `json:"guess"`
		NextChar string `json:"nextchar"`
	}

	type Result struct {
		Id            int64  `json:"id"`
		NumGuesses    int64  `json:"numguesses"`
		CurrentString string `json:"currentstring"`
		Correct       bool   `json:"correct"`
	}

	type Game struct {
		NumGuesses int64
		Word       string
		MaskedWord string
		Name       string
	}
	*/

	// Get the variables we care about
	id := ctx.RequireInt("id")
	guess, _ := ctx.FindString("guess")
	var nextchar string
	//FIXME
	//nextchar = ?

	// Sanitize inputs
	var game Game
	var ok bool
	if game, ok = games[id]; !ok {
		ctx.Error = jas.NewRequestError("Given game ID not found")
		panic(nil)
	}
	if len(nextchar) != 1 {
		//FIXME
	}

	// Increment the number of guesses
	//FIXME
	//game.?

	// Create our base result value
	result := Result{
		Id: id,
		NumGuesses: game.NumGuesses,
		Correct: false,
	}

	// See if the guess is correct; if so, set Correct to true,
	// set the returned word to the fully revealed value, delete the game from memory
	// and return the result via ctx.Data
	if len(guess) != 0 {
		if guess == game.Word {
			//FIXME
			return
		}
	}

	// otherwise, check the character that they gave us
	game.MaskedWord = UnhideByte([]byte(nextchar)[0], game.Word, game.MaskedWord)
	result.CurrentString = game.MaskedWord
	games[id] = game
	ctx.Data = result
	return
}
