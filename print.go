package main

import (
	"net/http"
	"os"
	"text/template"
)

func matchesPage(w http.ResponseWriter, r *http.Request) {

	t, err := template.New("matches").ParseFiles("matches.html")
	_ErrLogExit0(err)

	if r.Method == "POST" {
		switch submit := r.FormValue("submit"); submit {
		case "exit":
			f("exit matches")
			os.Exit(deInit())
		case "add":
			f("matches redir -> players")
			http.Redirect(w, r, "/players", http.StatusSeeOther)
		default:
			f("matches")
			err := t.ExecuteTemplate(w, "matches.html", players)
			_ErrLogExit0(err)
		}
	} else {
		f("matches")
		err := t.ExecuteTemplate(w, "matches.html", players)
		_ErrLogExit0(err)
	}
}

func playersPage(w http.ResponseWriter, r *http.Request) {

	t, err := template.New("players").ParseFiles("players.html")
	_ErrLogExit0(err)

	if r.Method == "POST" {
		switch submit := r.FormValue("submit"); submit {
		case "add":
			f("add player")
			if err := r.ParseForm(); err != nil {
				_ErrLogExit0(err)
			}
			who := r.FormValue("who")
			link := r.FormValue("link")
			userAdd(who, link)
			http.Redirect(w, r, "/players", http.StatusSeeOther)
		case "mathes":
			f("players redir -> matches")
			for _, p := range players {
				p.Matches = matchesFromOpenDotaAPI(p.SteamID32)
			}
			http.Redirect(w, r, "/matches", http.StatusSeeOther)
		case "update":
			for k := range players {
				//can do simultaneously
				updatePlayer(k)
			}
			http.Redirect(w, r, "/players", http.StatusSeeOther)
		case "exit":
			f("exit players")
			os.Exit(deInit())
		default:
			f("players")
			err = t.ExecuteTemplate(w, "players.html", players)
			_ErrLogExit0(err)
		}
	} else {
		f("players")
		err = t.ExecuteTemplate(w, "players.html", players)
		_ErrLogExit0(err)
	}
}
