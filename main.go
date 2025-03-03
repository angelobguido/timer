package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Structure for the expected JSON request
type IncomingRequest struct {
	ID      string          `json:"id"`
	Time    int             `json:"time"`
	Request json.RawMessage `json:"request"` // Preserve request JSON
}

// Structure for the "request" object inside the JSON
type ForwardRequest struct {
	URL     string            `json:"url"`
	Method  string            `json:"method"`
	Headers map[string]string `json:"headers"`
	Body    json.RawMessage   `json:"body,omitempty"` // Optional body
}

//var timers = make(map[string]time.Timer)

func main() {
	router := gin.Default()
	router.POST("/reset-timer", forwardRequest)

	router.Run("localhost:8080")
}

func forwardRequest(c *gin.Context) {
	// Parse incoming request
	var reqData IncomingRequest
	if err := c.ShouldBindJSON(&reqData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	fmt.Println("Received ID:", reqData.ID)
	fmt.Println("Received Time:", reqData.Time)

	// Parse the "request" field
	var forwardReq ForwardRequest
	if err := json.Unmarshal(reqData.Request, &forwardReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Create the request body if needed
	var reqBody io.Reader
	if len(forwardReq.Body) > 0 {
		reqBody = bytes.NewReader(forwardReq.Body)
	} else {
		reqBody = nil
	}

	// Create the HTTP request dynamically
	req, err := http.NewRequest(forwardReq.Method, forwardReq.URL, reqBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	// Add headers
	for key, value := range forwardReq.Headers {
		req.Header.Set(key, value)
	}

	// Perform the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to forward request"})
		return
	}

	resp.Body.Close()
}
