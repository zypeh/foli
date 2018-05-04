package main

import (
	"net/http"

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
	g := gin.Default()

	g.GET("/", toJSON)
	g.POST("/echo", queryJSON)
	// g.POST("/query", queryJSON)
	g.Run() // default localhost:8080
}

func toJSON(c *gin.Context) {
	// c.JSON only accept interface{}
	// https://github.com/gin-gonic/gin/issues/87
	//
	responses := []Query{
		{"title1", "description1", "filename1", "src1"},
		{"title2", "description2", "filename2", "src2"},
	}
	var respSlice []interface{} = make([]interface{}, len(responses))
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
