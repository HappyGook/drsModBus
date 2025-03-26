package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goburrow/modbus"
	"go.bug.st/serial"
	"log"
	"net/http"
	"time"
)

// REGISTER ID's
const (
	VOUT_SET          = 0x0020 //Output Voltage set
	CURVE_CV          = 0x00B1 // Constant voltage setting
	CURVE_FV          = 0x00B2 // Floating voltage setting
	CURVE_CC_TIMEOUT  = 0x00B5 // CC charge timeout setting
	CURVE_CV_TIMEOUT  = 0x00B6 // CV charge timeout setting
	CURVE_FV_TIMEOUT  = 0x00B7 // FV charge timeout setting
	BAT_UVP_SET       = 0x00D0 // BAT_LOW protect setting
	Force_BAT_UVP_SET = 0x00D1 // Force BAT_LOW protect setting
	BAUD_RATE         = 115200
)

// An array with register names of the meanwell drs
var registers = []uint16{
	VOUT_SET,
	CURVE_CV,
	CURVE_FV,
	CURVE_CC_TIMEOUT,
	CURVE_CV_TIMEOUT,
	CURVE_FV_TIMEOUT,
	BAT_UVP_SET,
	Force_BAT_UVP_SET,
}

// DRSClient A general struct with DRS features
type DRSClient struct {
	handler *modbus.RTUClientHandler
	client  modbus.Client
}

// NewDRSClient Establishing a connection with a DRS using a given port
func NewDRSClient(port string, baud int) (*DRSClient, error) {
	handler := modbus.NewRTUClientHandler(port)
	handler.BaudRate = baud
	handler.DataBits = 8
	handler.StopBits = 1
	handler.Parity = "N" // For NO parity
	handler.Timeout = 10 * time.Second
	handler.SlaveId = 131

	err := handler.Connect()
	if err != nil {
		log.Println("Could not connect to the device: ", err)
		return nil, err
	}

	client := modbus.NewClient(handler)
	return &DRSClient{handler: handler, client: client}, nil

}

// Methods of the Client struct

func (d *DRSClient) Close() {
	err := d.handler.Close()
	if err != nil {
		return
	}
}

func (d *DRSClient) ReadRegisters() ([]uint16, error) {
	results := make([]uint16, len(registers))

	for i, r := range registers {
		data, err := d.client.ReadHoldingRegisters(r, 1)
		if err != nil {
			return nil, err
		}
		results[i] = uint16(data[0])<<8 | uint16(data[1])
	}
	return results, nil
}

func (d *DRSClient) WriteRegisters(numbers []uint16) error {
	for i, r := range registers {
		_, err := d.client.WriteSingleRegister(r, numbers[i])
		if err != nil {
			return err
		}
	}
	return nil
}

// Handling functions
func handleSubmit(c *gin.Context) {
	port := c.Query("port") // Get port from query parameter
	if port == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Port is required"})
		log.Println("Port is required, but not received")
		return
	}

	var data struct {
		Values []uint16 `json:"values"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Println("Could not bind data: ", err)
		return
	}

	if len(data.Values) != len(registers) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect number of values"})
		log.Println("Incorrect number of values")
		return
	}

	drs, err := NewDRSClient(port, BAUD_RATE)
	if err != nil {
		log.Println("Failed to connect (when submitting)", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer drs.Close()

	err = drs.WriteRegisters(data.Values)
	if err != nil {
		log.Println("Failed to write (after submitting)", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data})
	fmt.Println("Data! :) ", data)

}

func handleRead(c *gin.Context) {
	port := c.Query("port") // Get port from query parameter
	if port == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Port is required"})
		return
	}
	log.Println("✅ Received request to read from port:", port)

	drs, err := NewDRSClient(port, BAUD_RATE)
	if err != nil {
		log.Println("Connection failed: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer drs.Close()
	log.Println("✅ Successfully connected to:", port)

	values, err := drs.ReadRegisters()
	if err != nil {
		log.Println("Failed to read the registers: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Println("✅ Successfully read registers:", values)
	c.JSON(http.StatusOK, gin.H{"registers": values})
}

func handleList(c *gin.Context) {
	ports, err := serial.GetPortsList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(ports) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No serial ports found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ports": ports})
}

func main() {
	r := gin.Default()

	r.GET("/api/list", handleList)
	r.POST("/api/submit", handleSubmit)
	r.GET("/api/read", handleRead)

	err := r.Run(":8080")
	if err != nil {
		return
	}

}
