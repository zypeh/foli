package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// has to be `type` to return JSON array
// https://github.com/gin-gonic/gin/issues/87
type response struct {
	Title      string `json:"title"`
	Descrption string `json:"description"`
	Filename   string `json:"filename"`
	Src        string `json:"src"`
}

func main() {
	g := gin.Default()

	g.GET("/", toJSON)
	g.Run() // default localhost:8080
}

func toJSON(c *gin.Context) {
	// c.JSON only accept interface{}
	// https://github.com/gin-gonic/gin/issues/87
	//
	responses := []response{
		{"title1", "description1", "filename1", "src1"},
		{"title2", "description2", "filename2", "src2"},
	}
	var respSlice []interface{} = make([]interface{}, len(responses))
	for i, resp := range responses {
		respSlice[i] = resp
	}

	fmt.Printf("%+v\n", respSlice)
	c.JSON(http.StatusOK, respSlice)
}
