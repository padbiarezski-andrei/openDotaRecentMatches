package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

const matchesJSONPPath = "./test/testRecentMatches.json"

type match struct {
	MatchID    int64 `json:"match_id"`
	PlayerSlot int   `json:"player_slot"`
	RadiantWin bool  `json:"radiant_win"`
	// Duration     int   `json:"duration"`
	// GameMode     int   `json:"game_mode"`
	// LobbyType    int   `json:"lobby_type"`
	// HeroID       int   `json:"hero_id"`
	// StartTime    int   `json:"start_time"`
	// Version      int   `json:"version"`
	// Kills        int   `json:"kills"`
	// Deaths       int   `json:"deaths"`
	// Assists      int   `json:"assists"`
	// Skill        int   `json:"skill"`
	// XpPerMin     int   `json:"xp_per_min"`
	// GoldPerMin   int   `json:"gold_per_min"`
	// HeroDamage   int   `json:"hero_damage"`
	// TowerDamage  int   `json:"tower_damage"`
	// HeroHealing  int   `json:"hero_healing"`
	// LastHits     int   `json:"last_hits"`
	// Lane         int   `json:"lane"`
	// LaneRole     int   `json:"lane_role"`
	// IsRoaming    bool  `json:"is_roaming"`
	// Cluster      int   `json:"cluster"`
	// LeaverStatus int   `json:"leaver_status"`
	// PartySize    int   `json:"party_size"`
}

func matchesFromOpenDotaAPI(steam32 string) []*match {
	f("get matches for player fron opendota API")
	r, err := http.Get("https://api.opendota.com/api/players/" + steam32 + "/recentMatches")
	_ErrLogExit0(err)

	var res []*match

	err = json.NewDecoder(r.Body).Decode(&res)
	_ErrLogExit0(err)

	for _, m := range res {
		if (m.RadiantWin && (m.PlayerSlot&128) == 0) || (!m.RadiantWin && (m.PlayerSlot&128) == 128) {
			m.RadiantWin = true
		} else {
			m.RadiantWin = false
		}
	}
	return res
}

func matchesLoadJSONFromFile() {
	data, err := os.ReadFile(matchesJSONPPath)
	_ErrLogExit0(err)
	for _, p := range players {
		//tmp := make([]match, 0)
		err = json.Unmarshal(data, &p.Matches)
		_ErrLogExit0(err)
		//fmt.Println(tmp)
		//fmt.Println()
		for _, m := range p.Matches {
			if (m.RadiantWin && (m.PlayerSlot&128) == 0) || (!m.RadiantWin && (m.PlayerSlot&128) == 128) {
				m.RadiantWin = true
			} else {
				m.RadiantWin = false
			}
		}

	}
	// players["77777777777777777"].Matches[0].MatchID = 1234567
	for _, p := range players {
		fmt.Println(p.Matches)
		fmt.Println()
	}
	f("matches loaded")
}
