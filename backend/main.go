package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Input struct {
	Number []string `json:"numbers"`
}

func main() {
	r := gin.Default()

	r.POST("/api/submit", handleSubmit)

	err := r.Run(":8080")
	if err != nil {
		return
	}

}

func handleSubmit(c *gin.Context) {
	var data Input

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data})
	fmt.Println("Data! :) ", data)

}
