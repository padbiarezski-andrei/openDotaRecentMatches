package main

import (
	"fmt"
	"net/http"
)

var players map[string]*player
var matches map[string][]*match

func deInit() int {
	defer writePlayersToFile()
	return 0
}
func init() {

	f("init")

	players = make(map[string]*player)
	matches = make(map[string][]*match)
	loadPlayersFromFile()
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
