package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func HttpAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		fmt.Fprint(w, "Method Error")
		return
	}
	path := strings.Split(r.URL.Path, "/")
	/* Save game data */
	if path[len(path)-1] == "save" {
		err := ioutil.WriteFile(r.PostFormValue("name"), []byte(r.PostFormValue("data")), os.ModePerm)
		if err != nil {
			fmt.Fprint(w, "ERROR")
			log.Printf("Failed while writing '%s'.", r.PostFormValue("name"))
		} else {
			fmt.Fprint(w, "OK")
		}
	}
}

func HttpGame(w http.ResponseWriter, r *http.Request) {
	scene := r.URL.Query().Get("scene")
	_html, _ := ioutil.ReadFile("game.html")
	_game_data, _ := ioutil.ReadFile(fmt.Sprintf("%s.gws", scene))
	/* 替换game_data */
	_game_data = []byte(strings.ReplaceAll(string(_game_data), "\"", "\\\"")) // " -> \"
	game_data_list := strings.Split(string(_game_data), "\n")
	game_data := "[\n"
	i := 0
	for {
		if i == len(game_data_list) {
			game_data = game_data[:len(game_data)-1] //delete ','
			break
		}
		if game_data_list[i] == "" {
			i++
			continue
		}
		game_data = fmt.Sprintf("%s\n\"%s\",", game_data, game_data_list[i])
		i++
	}
	game_data = fmt.Sprintf("%s\n%s", game_data, "]")
	/* game_data: [
		"x",
		"x"
	] */
	html := strings.ReplaceAll(string(_html), "{game_data}", game_data)
	fmt.Fprint(w, html)
}

func HttpResource(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	data, _ := ioutil.ReadFile(fmt.Sprintf("resource/%s", path[len(path)-1]))
	fmt.Fprint(w, string(data))
	w.Header().Set("Content-Type", http.DetectContentType(data))
}

func main() {
	http.HandleFunc("/game", HttpGame)
	http.HandleFunc("/api/", HttpAPI)
	http.HandleFunc("/resource/", HttpResource)
	http.ListenAndServe(":5000", nil)
}