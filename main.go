package main

import (
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
	ensureEnv("API")

	g := gin.Default()
	g.GET("/", toJSON)
	g.POST("/echo", queryJSON)
	// g.POST("/query", queryJSON)
	g.Run() // default localhost:8080
}

func ensureEnv(key string) {
	val, ok := os.LookupEnv(key)
	if !ok {
		// Early returning, return 1 to end this process
		// https://stackoverflow.com/questions/33885235/should-a-go-package-ever-use-log-fatal-and-when
		log.Fatalf("%s not set, exiting ...\nShould provide the Behance api key or client id in order to query images.", key)
	} else {
		fmt.Printf("Using \"%s\" as API key / client ID...\n\n", val)
	}
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
