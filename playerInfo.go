package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"unicode"
)

const playersJSONPPath = "./test/players.json"
const keyAPI = ""

var needUpdateJSON bool

type player struct {
	Who         string   `json:"who"`
	PersonaName string   `json:"persona_name"`
	SteamID32   string   `json:"steam_id_32"`
	SteamID64   string   `json:"steam_id_64"`
	LastUpdate  string   `json:"last_update"`
	AvatarFull  string   `json:"avatarfull"`
	Matches     []*match `json:"matches"`
}

func playersLoadJSONFromFile() {
	data, err := os.ReadFile(playersJSONPPath)
	_ErrLogExit0(err)

	err = json.Unmarshal(data, &players)
	fmt.Println(err)
	_ErrLogExit0(err)
	f("players loaded")
}

func playersWriteJSONToFile() {
	if needUpdateJSON {
		f("write players")
		jsonString, _ := json.Marshal(&players)
		os.WriteFile(playersJSONPPath, jsonString, os.ModePerm)
	}
}

func stringToSteamID(str string) (string, string, error) {
	if i := strings.Index(str, "https://steamcommunity.com/profiles/"); i != -1 {
		tmp := str[len("https://steamcommunity.com/profiles/"):]
		if !unicode.IsDigit(rune(tmp[len(tmp)-1])) {
			tmp = tmp[:len(tmp)-1]
		}
		n, err := strconv.ParseInt(tmp, 10, 64)
		_ErrLogExit0(err)
		return tmp, strconv.FormatInt(n-(int64(76561197960265728)), 10), nil
	}
	if i := strings.Index(str, "https://steamcommunity.com/id/"); i != -1 {
		tmp := str[len("https://steamcommunity.com/id/"):]
		if !unicode.IsDigit(rune(tmp[len(tmp)-1])) {
			tmp = tmp[:len(tmp)-1]
		}
		r, err := http.Get("http://api.steampowered.com/ISteamUser/ResolveVanityURL/v0001/?key=" + keyAPI + "&format=json&vanityurl=" + tmp)
		_ErrLogExit0(err)

		res := struct {
			Response struct {
				Steamid string `json:"steamid"`
				Success int    `json:"success"`
			} `json:"response"`
		}{}

		json.NewDecoder(r.Body).Decode(&res)
		if res.Response.Success == 42 {
			return "", "", fmt.Errorf("ERROR! stringToSteamID cannot extract steam id GET PARSE")
		}
		n, err := strconv.ParseInt(res.Response.Steamid, 10, 64)
		_ErrLogExit0(err)
		return res.Response.Steamid, strconv.FormatInt(n-(int64(76561197960265728)), 10), nil
	}
	if i := strings.LastIndex(str, "https://www.dotabuff.com/players/"); i != -1 {
		n, err := strconv.ParseInt(str[len("https://www.dotabuff.com/players/"):], 10, 64)
		_ErrLogExit0(err)
		return strconv.FormatInt((int64(76561197960265728) + n), 10), str[len("https://www.dotabuff.com/players/"):], nil
	}
	return "", "", fmt.Errorf("ERROR! stringToSteamID cannot extract steam id")
}

func userAdd(who, strInfo string) {
	// Who         string    `json:"who"`
	// PersonaName string    `json:"persona_name"`
	// AccountID   string    `json:"account_id"`
	// SteamID     string    `json:"steam_id"`
	// LastUpdate  time.Time `json:"last_update"`
	// AvatarFull  string    `json:"avatarfull"`
	f("adding " + who + " " + strInfo)
	steamid64, steamid32, _ := stringToSteamID(strInfo)

	if _, ok := players[steamid64]; ok {
		f("alredy present " + who + " " + strInfo)
	} else {
		r, err := http.Get("https://api.steampowered.com/ISteamUser/GetPlayerSummaries/v2/?key=" + keyAPI + "&format=json&steamids=" + steamid64)
		_ErrLogExit0(err)

		res := struct {
			Response struct {
				Players []struct {
					// Steamid                  string `json:"steamid"`
					// Communityvisibilitystate int    `json:"communityvisibilitystate"`
					// Profilestate             int    `json:"profilestate"`
					Personaname string `json:"personaname"`
					// Commentpermission        int    `json:"commentpermission"`
					// Profileurl               string `json:"profileurl"`
					// Avatar                   string `json:"avatar"`
					// Avatarmedium             string `json:"avatarmedium"`
					Avatarfull string `json:"avatarfull"`
					// Avatarhash               string `json:"avatarhash"`
					// Personastate             int    `json:"personastate"`
					// Realname                 string `json:"realname"`
					// Primaryclanid            string `json:"primaryclanid"`
					// Timecreated              int    `json:"timecreated"`
					// Personastateflags        int    `json:"personastateflags"`
					// Loccountrycode           string `json:"loccountrycode"`
					// Locstatecode             string `json:"locstatecode"`
					// Loccityid                int    `json:"loccityid"`
				} `json:"players"`
			} `json:"response"`
		}{}

		json.NewDecoder(r.Body).Decode(&res)
		if len(res.Response.Players) == 0 {
			_ErrLogExit0(fmt.Errorf("ERROR! stringToSteamID cannot extract steam id GET PARSE"))
		}

		players[steamid64] = &player{
			Who:         who,
			PersonaName: res.Response.Players[0].Personaname,
			SteamID32:   steamid32,
			SteamID64:   steamid64,
			LastUpdate:  timeStamp(),
			AvatarFull:  res.Response.Players[0].Avatarfull,
		}
		f("added " + who + " " + strInfo)
		f("get matches for player fron opendota API")
		rm, err := http.Get("https://api.opendota.com/api/players/" + steamid32 + "/recentMatches")
		_ErrLogExit0(err)

		var matchs []*match

		err = json.NewDecoder(rm.Body).Decode(&matchs)
		_ErrLogExit0(err)

		for _, m := range matchs {
			if (m.RadiantWin && (m.PlayerSlot&128) == 0) || (!m.RadiantWin && (m.PlayerSlot&128) == 128) {
				m.RadiantWin = true
			} else {
				m.RadiantWin = false
			}
		}
		players[steamid64].Matches = matchs
		needUpdateJSON = true
	}
}

func updatePlayer(steamid64 string) error {
	if _, ok := players[steamid64]; ok {
		f("updating " + steamid64)
		steamid32 := players[steamid64].SteamID32
		who := players[steamid64].Who
		r, err := http.Get("https://api.steampowered.com/ISteamUser/GetPlayerSummaries/v2/?key=" + keyAPI + "&format=json&steamids=" + steamid64)
		_ErrLogExit0(err)

		res := struct {
			Response struct {
				Players []struct {
					// Steamid                  string `json:"steamid"`
					// Communityvisibilitystate int    `json:"communityvisibilitystate"`
					// Profilestate             int    `json:"profilestate"`
					Personaname string `json:"personaname"`
					// Commentpermission        int    `json:"commentpermission"`
					// Profileurl               string `json:"profileurl"`
					// Avatar                   string `json:"avatar"`
					// Avatarmedium             string `json:"avatarmedium"`
					Avatarfull string `json:"avatarfull"`
					// Avatarhash               string `json:"avatarhash"`
					// Personastate             int    `json:"personastate"`
					// Realname                 string `json:"realname"`
					// Primaryclanid            string `json:"primaryclanid"`
					// Timecreated              int    `json:"timecreated"`
					// Personastateflags        int    `json:"personastateflags"`
					// Loccountrycode           string `json:"loccountrycode"`
					// Locstatecode             string `json:"locstatecode"`
					// Loccityid                int    `json:"loccityid"`
				} `json:"players"`
			} `json:"response"`
		}{}

		json.NewDecoder(r.Body).Decode(&res)
		if len(res.Response.Players) == 0 {
			_ErrLogExit0(fmt.Errorf("ERROR! stringToSteamID cannot extract steam id GET PARSE"))
		}

		players[steamid64] = &player{
			Who:         who,
			PersonaName: res.Response.Players[0].Personaname,
			SteamID32:   steamid32,
			SteamID64:   steamid64,
			LastUpdate:  timeStamp(),
			AvatarFull:  res.Response.Players[0].Avatarfull,
		}
		f("update " + steamid64)

		f("get matches for player fron opendota API")
		rm, err := http.Get("https://api.opendota.com/api/players/" + steamid32 + "/recentMatches")
		_ErrLogExit0(err)

		var matchs []*match

		err = json.NewDecoder(rm.Body).Decode(&matchs)
		_ErrLogExit0(err)

		for _, m := range matchs {
			if (m.RadiantWin && (m.PlayerSlot&128) == 0) || (!m.RadiantWin && (m.PlayerSlot&128) == 128) {
				m.RadiantWin = true
			} else {
				m.RadiantWin = false
			}
		}
		players[steamid64].Matches = matchs
		needUpdateJSON = true
	} else {
		userAdd("who", steamid64)
	}
	return nil
}
