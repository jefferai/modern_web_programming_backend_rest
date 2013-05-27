package main

import (
	"github.com/coocood/jas"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// STARTRECEIVESTRUCT OMIT
// This will not actually be used, but it will be our reference
// for what we expect the input data to look like
type ReceivedData struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	NewGame  bool   `json:"newgame"`
	Guess    string `json:"guess"`
	NextChar string `json:"nextchar"`
}
// ENDRECEIVESTRUCT OMIT

func sigintCatcher() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT)
	<-ch
	os.Exit(0)
}

func main() {
	// Catch signals
	go sigintCatcher()

	router := jas.NewRouter(new(Wordgame))
	router.RequestErrorLogger = router.InternalErrorLogger

	log.Println("Starting serving on paths:\n" + router.HandledPaths(true))
	http.Handle(router.BasePath, router)
	err := http.ListenAndServe("127.0.0.1:8888", nil)
	if err != nil {
		panic(err)
	}
}
