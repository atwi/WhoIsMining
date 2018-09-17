package main

import (
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
	"encoding/json"
)

func getMiners() map[string][]string {
	return map[string][]string{
		"jsecoin.com":             {`((.*?)load\.jsecoin\.com\/(.*?))`, `((.*?)server\.jsecoin\.com\/(.*?))`},
		"coin-hive.com":           {`((.*?)coin-hive\.com/lib(.*?))`, `((.*?)coin-hive\.com/proxy(.*?))`, `((.*?)coin-hive\.com/captcha(\.*?))`},
		"edgeno.de":               {`((.*?)edgeno\.de/(.*?))`},
		"reasedoper.pw":           {`((.*?)reasedoper\.pw/(.*?))`},
		"mataharirama.xyz":        {`((.*?)mataharirama\.xyz/(.*?))`},
		"listat.biz":              {`((.*?)listat\.biz/(.*?))`},
		"lmodr.biz":               {`((.*?)lmodr\.biz/(.*?))`},
		"minecrunch.co":           {`((.*?)minecrunch\.co/web/(.*?))`},
		"jyhfuqoh.info":           {`((.*?)jyhfuqoh\.info/(.*?))`},
		"coinhive.com":            {`((.*?)coinhive\.com/lib(.*?))`, `((.*?)coinhive\.com/proxy(.*?))`, `((.*?)coinhive\.com/captcha(.*?))`, `((.*?)coinhive\.js(.*?))`, `((.*?)coinhive\.min\.js(.*?))`},
		"minemytraffic.com":       {`((.*?)minemytraffic\.com/lib(.*?))`},
		"crypto-loot.com":         {`((.*?)crypto-loot\.com/lib(.*?))`, `((.*?)crypto-loot\.com/proxy(.*?))`},
		"2giga.link":              {`((.*?)2giga\.link/wproxy(.*?))`, `((.*?)2giga\.link/hive/lib/(.*?))`},
		"coin-have.com":           {`((.*?)coin-have\.com/c/(.*?))`},
		"ppoi.org":                {`((.*?)ppoi\.org/lib/(.*?))`, `((.*?)ppoi\.org/token/(.*?))`},
		"coinerra.com":            {`((.*?)coinerra\.com/lib/(.*?))`},
		"kisshentai.net":          {`((.*?)kisshentai\.net/Content/js/c-hive\.js((.*?)`},
		"miner.pr0gramm.com":      {`((.*?)miner\.pr0gramm\.com/xmr\.min\.js(.*?))`},
		"minero.pw":               {`((.*?)minero\.pw/miner\.min\.js(.*?))`},
		"kiwifarms.net":           {`((.*?)kiwifarms\.net/js/Jawsh/xmr/xmr\.min\.js(.*?))`},
		"coinblind.com":           {`((.*?)coinblind\.com/lib/(.*?))`},
		"joyreactor.cc":           {`((.*?)joyreactor\.cc/ws/ch/(.*?))`},
		"reactor.cc":              {`((.*?)anime\.reactor\.cc/js/ch/cryptonight\.wasm(.*?))`},
		"webmine.cz":              {`((.*?)webmine\.cz/miner(.*?))`},
		"monerominer.rocks":       {`((.*?)monerominer\.rocks/scripts/miner\.js(.*?))`, `((.*?)monerominer\.rocks/miner\.php(.*?))`},
		"kissdoujin.com":          {`((.*?)kissdoujin\.com/Content/js/c-hive\.js(.*?))`},
		"inwemo.com":              {`((.*?)inwemo\.com/inwemo\.min\.js(.*?))`},
		"coinnebula.com":          {`((.*?)coinnebula\.com/lib/(.*?))`},
		"afminer.com":             {`((.*?)afminer\.com/code/(.*?))`},
		"cloudcoins.co":           {`((.*?)cloudcoins\.co/javascript/(.*?))`},
		"coinlab.biz":             {`((.*?)coinlab\.biz/lib/coinlab\.js(.*?))`},
		"papoto.com":              {`((.*?)papoto\.com/lib/(.*?))`},
		"cookiescript.info":       {`((.*?)cookiescript\.info/libs/(.*?))`},
		"cookiescriptcdn.pro":     {`((.*?)cookiescriptcdn\.pro/libs/(.*?))`},
		"ad-miner.com":            {`((.*?)ad-miner\.com/lib/(.*?))`},
		"party-nngvitbizn.now.sh": {`((.*?)party-nngvitbizn\.now\.sh(.*?))`},
		"rocks.io":                {`((.*?)rocks\.io/assets/(.*?))`, `((.*?)rocks\.io/proxy/(.*?))`},
		"coinhive-manager.com":    {`((.*?)coinhive-manager\.com\/ch-manager\.js(.*?))`},
	}
}

var t *template.Template

// SearchResult represents the single search attempt done by the user and shown once in the user interface
type SearchResult struct {
	URL   string `json:"url"`
	Miner string `json:"miner"`
}

// Website represents the page to be scanned
type Website struct {
	URL  string
	Body []byte
}

// NewWebsite creates a new website from url and download index page content
func NewWebsite(url string) (*Website, error) {
	url = strings.ToLower(url)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return NewEmptyWebsite(url), errors.New("http.Get: Invalid address")
	}
	req.Header.Set("User-Agent", "whoismining.com Bot/1.0")
	res, err := client.Do(req)
	if err != nil {
		return NewEmptyWebsite(url), errors.New("http.Get: Invalid address")
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return NewEmptyWebsite(url), errors.New("http.Get: Unable to fetch page content")
	}
	resBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return NewEmptyWebsite(url), err
	}
	return &Website{url, resBytes}, nil
}

// NewEmptyWebsite returns a empty instance of the Website struct
// called when is impossible to reach the host
func NewEmptyWebsite(url string) *Website {
	return &Website{url, make([]byte, 0)}
}

// IsMining search the miners in the index page of the website
// TODO(0x13a): add recursive search to other js scripts in the page as right now it just scan the index of the website
func (w *Website) IsMining(minersList map[string][]string) (bool, string) {
	for provider, scripts := range minersList {
		for _, s := range scripts {
			match, _ := regexp.MatchString(s, string(w.Body))
			if match == false {
				continue
			}
			return true, provider
		}
	}
	return false, ""
}

// NewSearchResult creates a new search result
func NewSearchResult(url string, miner string) *SearchResult {
	return &SearchResult{url, miner}
}

// LogSearch
func LogSearch(r *SearchResult) {
	f, err := os.OpenFile("./search.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer f.Close()
	if err == nil {
		fmt.Fprintf(f, "%d,%s,%s\n", time.Now().Unix(), r.URL, r.Miner)
	}
}

func loadTemplate() {
	t = template.Must(template.ParseFiles("search.html"))
}

// MethodNotSupported returns to the home page when a method is not supported by the web server
func MethodNotSupported(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "http://www.whoismining.com", 301)
}

// Index shows the index page
func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	err := t.Execute(w, nil)
	if err != nil {
		fmt.Fprintf(w, "Oops, something went wrong. Try to reload the page")
	}
}

// Check is the http handler that scan the website
func Check(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	url := r.FormValue("q")
	match, err := regexp.Match(`^(?:f|ht)tps?://`, []byte(url))
	if match == false {
		url = "http://" + url
	}
	ws, err := NewWebsite(url)
	if err != nil {
		json.NewEncoder(w).Encode(NewSearchResult("", ""))
		return
	}
	isMining, miner := ws.IsMining(getMiners())
	if isMining {
		LogSearch(NewSearchResult(ws.URL, miner))
		json.NewEncoder(w).Encode(NewSearchResult(ws.URL, miner))
		return
	}
	json.NewEncoder(w).Encode(NewSearchResult(ws.URL, ""))
}

// HandleSearch is the http.Handler
// it will handle the request based on the method and query param "q"
// it will either check for a website or return the landing page
func HandleSearch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		MethodNotSupported(w, r)
		return
	}
	if r.FormValue("q") != "" {
		Check(w, r)
		return
	}
	Index(w, r)
}

func main() {
	loadTemplate()
	http.HandleFunc("/", HandleSearch)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
