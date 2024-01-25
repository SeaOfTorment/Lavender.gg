package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"text/template"
)

//Types
type Handler int

type stats struct {
	User string `json:"name"`
	Lvl int `json:"summonerLevel"`
	Pfp int `json:"profileIconId"`
}

//Globals
var tpl *template.Template

//Code
func init() {
	tpl = template.Must(template.ParseGlob("Templates/*"))
}

func main() {
	fmt.Println("Lavender.gg ONLINE!")
	var handler Handler
	http.ListenAndServe(":80", handler)
}

func (handler Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(1000000)
	if r.URL.Path == "/" {
		//fmt.Println(r.PostForm)
		if len(r.PostForm) == 0 {
			tpl.ExecuteTemplate(w, "main.gohtml", nil)
		} else if len(r.PostForm) == 2 {
			searchMux(r.PostForm, w)
		} else {
			tpl.ExecuteTemplate(w, "main.gohtml", nil)
		}
	} else if r.URL.Path == "/favicon.ico" {
		http.ServeFile(w, r, "Files/favicon.svg")
	} else {
		fmt.Fprintln(w, "404:1")
	}
}

func searchMux(pf url.Values, w http.ResponseWriter) {
	if len(pf["Username"][0]) != 0 {
		resp, err := http.Get("https://na1.api.riotgames.com/lol/summoner/v4/summoners/by-name/"+pf["Username"][0]+"?api_key="+os.Getenv("api"))
		if err != nil {
  		log.Fatalln(err)
		}
		responseData, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Fatal(err)
    }
		loaded_stats := stats{}
    r := []byte(responseData)
		err = json.Unmarshal(r, &loaded_stats)
		fmt.Println("Query: " + loaded_stats.User)
		tpl.ExecuteTemplate(w, "main.gohtml", loaded_stats)
		return
	}
}