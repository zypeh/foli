package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/asdine/storm"
	"github.com/gin-gonic/gin"
)

// has to be `type` to return JSON array
// https://github.com/gin-gonic/gin/issues/87
//
// Using
//
type Query struct {
	Title      string `json:"title"`
	Descrption string `json:"description"`
	Filename   string `json:"filename"`
	Src        string `json:"src"`
}

// JSON parsing and accessing
// https://github.com/astaxie/build-web-application-with-golang/blob/master/zh/07.2.md
type CreativesSlice struct {
	Creatives []Creative `json:"creatives_to_follow"`
}

type Creative struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

type UserProjectsSlice struct {
	Projects []UserProject `json:"projects"`
}

type UserProject struct {
	ID int `json:"id"`
}

type Project struct {
	Project ProjectParsed `json:"project"`
}

type ProjectParsed struct {
	Title       string                 `json:"name"`
	Description string                 `json:"description"`
	Src         map[string]interface{} `json:"covers"`
}

type Data struct {
	ID          int    `storm:"id,increment" json:"id"`
	Title       string `storm:"index" json:"title"`
	Description string `json:"description"`
	Filename    string `storm:"index" json:"filename"`
	Src         string `storm:"index" json:"src"`
}

type Env struct {
	db *storm.DB
}

func main() {
	var api = ensureEnv("API")

	db, err := storm.Open(filepath.Join(".", "foli.db"))
	if err != nil {
		log.Fatalf("%s\n", err)
	}
	defer db.Close()

	fetchItem(api, db)

	g := gin.Default()
	env := &Env{db: db}

	g.GET("/", env.queryAll)
	g.Run() // default localhost:8080
}

func ensureEnv(key string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		// Early returning, return 1 to end this process
		// https://stackoverflow.com/questions/33885235/should-a-go-package-ever-use-log-fatal-and-when
		log.Fatalf("%s not set, exiting ...\nShould provide the Behance api key or client id in order to query images.", key)
	} else {
		fmt.Printf("Using \"%s\" as API key / client ID...\n\n", val)
	}
	return val
}

func toJSON(c *gin.Context) {
	// c.JSON only accept interface{}
	// https://github.com/gin-gonic/gin/issues/87
	//
	responses := []Query{
		{"title1", "description1", "filename1", "src1"},
		{"title2", "description2", "filename2", "src2"},
	}
	var respSlice = make([]interface{}, len(responses))
	for i, resp := range responses {
		respSlice[i] = resp
	}
	c.JSON(http.StatusOK, respSlice)
}

// Dump all the entries in DB
func (e *Env) queryAll(c *gin.Context) {
	var respJSON []Data
	e.db.All(&respJSON)
	c.JSON(http.StatusOK, respJSON)
}

// Use endpoint /v2/creativestofollow to fetch a list of creatives to follow (user). [get 10]
// Use endpoint /v2/users/:username to fetch a list of projects created by user.
// Use endpoint /v2/projects/:id to fetch the cover and description needed.
// And, it accepts a parameter to do pagination.
// https://www.behance.net/dev/api/endpoints/9
func fetchItem(apiKey string, db *storm.DB) {
	fetch := func(url string, page int, dest interface{}) error {
		urlWithPage := fmt.Sprintf("%s?page=%d&client_id=%s", url, page, apiKey)
		response, err := http.Get(urlWithPage)
		if err != nil {
			log.Printf("%s\n", err)
			return err
		}
		defer response.Body.Close() // resource management
		json.NewDecoder(response.Body).Decode(dest)
		return nil
	}

	var userList CreativesSlice
	var projectList UserProjectsSlice
	var resource Project
	fetch("https://api.behance.net/v2/creativestofollow", 1, &userList)
	// fmt.Printf("%s\n", result.Creatives[1].Images["115"])
	fetch(fmt.Sprintf("https://api.behance.net/v2/users/%s/projects", userList.Creatives[0].Username), 1, &projectList)
	fetch(fmt.Sprintf("https://api.behance.net/v2/projects/%d", projectList.Projects[0].ID), 1, &resource)

	dataToDB := Data{
		Title:       resource.Project.Title,
		Description: resource.Project.Description,
		Filename:    getFilename(resource.Project.Src["original"].(string)),
		Src:         resource.Project.Src["original"].(string),
	}

	// TODO: buffer queue
	go fetchImages(dataToDB.Src)
	db.Save(&dataToDB)
}

func fetchImages(src string) error {
	client := http.Client{
		Transport: &http.Transport{
			Dial: func(network, addr string) (net.Conn, error) {
				return net.DialTimeout(network, addr, 3*time.Second)
			},
		},
	}

	resp, err := client.Get(src)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body) // reads until EOF, for byte[]
	if err != nil {
		return err
	}

	// Path setup
	path := filepath.Join(".", "images")
	os.MkdirAll(path, os.ModePerm)

	filename := getFilename(src)
	// saves to fs
	ioutil.WriteFile(filepath.Join(path, filename), b, 0644)
	return nil
}

func getFilename(src string) string {
	url := strings.Split(src, "/")
	if len(url) < 1 {
		log.Fatalf("invalid url %s\n", src)
	}
	return url[len(url)-1]
}
