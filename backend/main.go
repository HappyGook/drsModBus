package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goburrow/modbus"
	"log"
	"net/http"
)

// REGISTER ID's ----------------
const (
	VOUT_SET          = 0x0020 //Output Voltage set
	CURVE_CV          = 0x00B1 // Constant voltage setting
	CURVE_FV          = 0x00B2 // Floating voltage setting
	CURVE_CC_TIMEOUT  = 0x00B5 // CC charge timeout setting
	CURVE_CV_TIMEOUT  = 0x00B6 // CV charge timeout setting
	CURVE_FV_TIMEOUT  = 0x00B7 // FV charge timeout setting
	BAT_UVP_SET       = 0x00D0 // BAT_LOW protect setting
	Force_BAT_UVP_SET = 0x00D1 // Force BAT_LOW protect setting

	IP   = "192.168.1.100" //PLACEHOLDERS!!!
	PORT = "502"
)

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

type DRSClient struct {
	handler *modbus.TCPClientHandler
	client  modbus.Client
}

func NewDRSClient(ip, port string) (*DRSClient, error) {
	handler := modbus.NewTCPClientHandler(fmt.Sprintf("%s:%s", ip, port))

	err := handler.Connect()
	if err != nil {
		log.Println("Could not connect to the device: ", err)
		return nil, err
	}

	client := modbus.NewClient(handler)
	return &DRSClient{handler: handler, client: client}, nil

}

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

func handleSubmit(c *gin.Context) {
	ip := IP
	port := PORT
	var data struct {
		Values []uint16 `json:"values"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(data.Values) != len(registers) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect number of values"})
		return
	}

	drs, err := NewDRSClient(ip, port)
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
	ip := IP
	port := PORT

	drs, err := NewDRSClient(ip, port)
	if err != nil {
		log.Println("Connection failed: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer drs.Close()

	values, err := drs.ReadRegisters()
	if err != nil {
		log.Println("Failed to read the registers: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"registers": values})
}

func main() {
	r := gin.Default()

	r.POST("/api/submit", handleSubmit)
	r.GET("/api/read", handleRead)

	err := r.Run(":8080")
	if err != nil {
		return
	}

}

/*
func drsCom() {
	ip := "192.168.1.100" // PLACEHOLDER !!!!
	port := "502"

	handler := modbus.NewTCPClientHandler(fmt.Sprintf("%s:%s", ip, port))
	handler.Timeout = 10 * time.Second

	err := handler.Connect()
	if err != nil {
		log.Fatal("Could not connect to the device: ", err)
	}
	defer func(handler *modbus.TCPClientHandler) {
		err := handler.Close()
		if err != nil {

		}
	}(handler)

	client := modbus.NewClient(handler)

	//Read Register Values to show them to the User
	registers := map[string]uint16{
		"Output Voltage Set":        VOUT_SET,
		"Constant Voltage Setting":  CURVE_CV,
		"Floating Voltage Setting":  CURVE_FV,
		"CC Charge Timeout Setting": CURVE_CC_TIMEOUT,
		"CV Charge Timeout Setting": CURVE_CV_TIMEOUT,
		"FV Charge Timeout Setting": CURVE_FV_TIMEOUT,
		"BAT_LOW Protect Setting":   BAT_UVP_SET,
		"Force BAT_LOW Protect":     Force_BAT_UVP_SET,
	}

	registerValues := make(map[string]uint16)

	for name, address := range registers {
		results, err := client.ReadHoldingRegisters(address, 1)
		if err != nil {
			log.Printf("Read of the Register (%s) failed: %v", name, err)
			continue
		}
		value := uint16(results[0])<<8 | uint16(results[1])
		registerValues[name] = value
	}

	for name, value := range registerValues {
		fmt.Printf("%s: %d\n", name, value)
	}

}
*/
