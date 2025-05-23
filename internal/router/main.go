package router

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	healthcheck "github.com/RaMin0/gin-health-check"
	brotli "github.com/anargu/gin-brotli"
	"github.com/bilte-co/bilte/internal/domain"
	"github.com/bilte-co/bilte/internal/templates"
	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
	stats "github.com/semihalev/gin-stats"
	"golang.org/x/time/rate"
)

// It keeps a list of clients those are currently attached
// and broadcasting events to those clients.
type Event struct {
	// Events are pushed to this channel by the main events-gathering routine
	Message chan string

	// New client connections
	NewClients chan chan string

	// Closed client connections
	ClosedClients chan chan string

	// Total client connections
	TotalClients map[chan string]bool
}

// New event messages are broadcast to all registered client connection channels
type ClientChan chan string

var limiter = rate.NewLimiter(50, 200)

func rateLimiter(c *gin.Context) {
	if !limiter.Allow() {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "too many requests"})
		c.Abort()
		return
	}
	c.Next()
}

func HeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		c.Writer.Header().Set("Transfer-Encoding", "chunked")
		c.Next()
	}
}

// Initialize event and Start procnteessing requests
func NewServer() (event *Event) {
	event = &Event{
		Message:       make(chan string),
		NewClients:    make(chan chan string),
		ClosedClients: make(chan chan string),
		TotalClients:  make(map[chan string]bool),
	}

	go event.listen()

	return
}

// It Listens all incoming requests from clients.
// Handles addition and removal of clients and broadcast messages to clients.
func (stream *Event) listen() {
	for {
		select {
		// Add new available client
		case client := <-stream.NewClients:
			stream.TotalClients[client] = true
			// publish a message
			stream.Message <- fmt.Sprint("Client added")
			log.Printf("Client added. %d registered clients", len(stream.TotalClients))

		// Remove closed client
		case client := <-stream.ClosedClients:
			delete(stream.TotalClients, client)
			close(client)
			stream.Message <- fmt.Sprint("Client removed")
			log.Printf("Removed client. %d registered clients", len(stream.TotalClients))

		// Broadcast message to client
		case eventMsg := <-stream.Message:
			for clientMessageChan := range stream.TotalClients {
				select {
				case clientMessageChan <- eventMsg:
					// Message sent successfully
					log.Printf("Message sent to client")
				default:
					// Failed to send, dropping message
					log.Printf("Failed to send message to client")
				}
			}
		}
	}
}

func (stream *Event) serveHTTP() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Initialize client channel
		clientChan := make(ClientChan)

		// Send new connection to event server
		stream.NewClients <- clientChan

		go func() {
			<-c.Request.Context().Done()

			// Drain client channel so that it does not block. Server may keep sending messages to this channel
			for range clientChan {
				fmt.Println("Draining client channel")
			}
			// Send closed connection to event server
			stream.ClosedClients <- clientChan
		}()

		c.Set("clientChan", clientChan)

		c.Next()
	}
}

func NewRouter(r *gin.Engine, production *bool, nc *nats.Conn) *gin.Engine {
	r.Use(brotli.Brotli(brotli.DefaultCompression))
	r.Use(healthcheck.Default())
	r.Use(stats.RequestStats())
	r.Use(rateLimiter)

	stream := NewServer()

	go func() {
		for {
			time.Sleep(time.Second * 10)
			now := time.Now().Format(time.RFC3339)
			currentTime := fmt.Sprintf("The Current Time Is %v", now)

			// Send current time to clients message channel
			fmt.Println(currentTime)
			stream.Message <- currentTime
		}
	}()

	r.GET("/stats", func(c *gin.Context) {
		c.JSON(http.StatusOK, stats.Report())
	})

	r.GET("/", func(c *gin.Context) {
		title := "bilte co"
		description := "Strategy-led software engineering and consulting—delivering impactful results in high-stakes domains."
		c.HTML(http.StatusOK, "", templates.Home(production, &title, &description))
	})

	r.GET("/cv", func(c *gin.Context) {
		var Info domain.Resume
		var Projects domain.Projects
		// we are going to get the data in data/resume.json and data/projects.json
		infoFile, err := os.Open("data/resume.json")
		if err != nil {
			c.String(http.StatusInternalServerError, "Error opening info file: "+err.Error())
			return
		}
		defer infoFile.Close()

		infoBytes, err := io.ReadAll(infoFile)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error reading info file: "+err.Error())
			return
		}
		err = json.Unmarshal(infoBytes, &Info)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error unmarshalling info file: "+err.Error())
			return
		}

		projectsFile, err := os.Open("data/projects.json")
		if err != nil {
			c.String(http.StatusInternalServerError, "Error opening projects file: "+err.Error())
			return
		}
		defer projectsFile.Close()

		projectsBytes, err := io.ReadAll(projectsFile)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error reading projects file: "+err.Error())
			return
		}

		err = json.Unmarshal(projectsBytes, &Projects)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error unmarshalling projects file: "+err.Error())
			return
		}

		c.HTML(http.StatusOK, "", templates.CV(production, Info, Projects))
	})

	r.GET("/sse", HeadersMiddleware(), stream.serveHTTP(), func(c *gin.Context) {
		v, ok := c.Get("clientChan")
		if !ok {
			return
		}
		clientChan, ok := v.(ClientChan)
		if !ok {
			return
		}
		c.Stream(func(w io.Writer) bool {
			// Stream message to client from message channel
			if msg, ok := <-clientChan; ok {
				fmt.Println("Streaming message to client:", msg)
				c.SSEvent("message", msg)
				return true
			}
			return false
		})
	})

	return r
}
