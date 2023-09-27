package main

import (
	"fmt"
	"net/http"
)

var players map[string]*player  // slice of tracking players
var matches map[string][]*match // slice of tracking matches

func deInit() int {
	defer playersWriteJSONToFile()
	return 0
}
func init() {

	f("init")
	players = make(map[string]*player)
	matches = make(map[string][]*match)
	playersLoadJSONFromFile()
}

func main() {
	defer deInit()

	http.HandleFunc("/matches", matchesPage)
	http.HandleFunc("/players", playersPage)
	f("starting server at :8080")
	openInBrowser("http://localhost:8080/players")
	http.ListenAndServe(":8080", nil)
	fmt.Scan()
}
