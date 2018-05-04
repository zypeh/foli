package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// has to be `type` to return JSON array
// https://github.com/gin-gonic/gin/issues/87
type Query struct {
	Title      string `json:"title"`
	Descrption string `json:"description"`
	Filename   string `json:"filename"`
	Src        string `json:"src"`
}

func main() {
	var api = ensureEnv("API")

	fetchItem(api)

	g := gin.Default()
	g.GET("/", toJSON)
	g.POST("/echo", queryJSON)
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

func queryJSON(c *gin.Context) {
	var queryJSON Query

	// Early return
	if c.BindJSON(&queryJSON) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Error occurred when parsing your JSON query ! X( "})
		return
	}

	c.JSON(http.StatusOK, queryJSON)
}

// Using /v2/creativestofollow to fetch a list of creatives.
// And, it accepts a parameter to do pagination.
// https://www.behance.net/dev/api/endpoints/9
func fetchItem(apiKey string) {
	fetch := func(url string, page int, dest interface{}) error {
		urlWithPage := fmt.Sprintf("%s?page=%d&client_id=%s", url, page, apiKey)
		response, err := http.Get(urlWithPage)
		if err != nil {
			log.Printf("%s\n", err)
		}
		defer response.Body.Close() // resource management
		return json.NewDecoder(response.Body).Decode(dest)
	}

	url := "https://api.behance.net/v2/creativestofollow"
	var foobar map[string]interface{}
	err := fetch(url, 1, &foobar)
	if err != nil {
		log.Fatalf("%s\n", err)
	}
	fmt.Printf("%+v\n", foobar)
}
