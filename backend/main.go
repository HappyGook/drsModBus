package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gith
	"log"
	"net/http"
	"time"
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

func drsCom() {
	ip := "192.168.1.100" // PLACEHOLDER !!!!
	port := "502"

	handler := modbus.NewTCPClientHandler(fmt.Sprintf("%s:%s", ip, port))
	handler.Timeout = 10 * time.Second

	err := handler.Connect()
	if err != nil {
		log.Fatal("Connection error: ", err)
	}
	defer func(handler *modbus.TCPClientHandler) {
		err := handler.Close()
		if err != nil {

		}
	}(handler)

	client := modbus.NewClient(handler)

	results, err := client.ReadHoldingRegisters(0x00, 2)
	if err != nil {
		log.Fatal("Reading failed: ", err)
	}

	fmt.Printf("Register values: %v\n\n", results)

}
