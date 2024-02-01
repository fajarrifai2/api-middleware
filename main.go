package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Ticket struct {
	ID       int    `json:"id"`
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Quantity int    `json:"quantity" binding:"required,min=1"`
	Token    string `json:"token"`
}

var tickets []Ticket
var idCounter = 1

func main() {
	router := gin.Default()

	router.Use(middlewareLogger())
	router.Use(middlewareAuth())

	router.GET("/tickets", getTickets)
	router.POST("/tickets", createTicket)
	router.GET("/tickets/:id", getTicket)
	router.PUT("/tickets/:id", updateTicket)
	router.DELETE("/tickets/:id", deleteTicket)

	router.Run(":8080")
}

func getTickets(c *gin.Context) {
	c.JSON(http.StatusOK, tickets)
}

func createTicket(c *gin.Context) {
	var ticket Ticket
	if err := c.ShouldBindJSON(&ticket); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ticket.ID = idCounter
	idCounter++

	ticket.Token = "JKT-SBY" // Set token secara statis sesuai kebutuhan
	tickets = append(tickets, ticket)
	c.JSON(http.StatusCreated, ticket)
}

func getTicket(c *gin.Context) {
	id := c.Param("id")
	for _, ticket := range tickets {
		if strconv.Itoa(ticket.ID) == id {
			c.JSON(http.StatusOK, ticket)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
}

func updateTicket(c *gin.Context) {
	id := c.Param("id")
	var updatedTicket Ticket
	if err := c.ShouldBindJSON(&updatedTicket); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for i, ticket := range tickets {
		if strconv.Itoa(ticket.ID) == id {
			tickets[i].Name = updatedTicket.Name
			tickets[i].Email = updatedTicket.Email
			tickets[i].Quantity = updatedTicket.Quantity
			c.JSON(http.StatusOK, tickets[i])
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
}

func deleteTicket(c *gin.Context) {
	id := c.Param("id")
	for i, ticket := range tickets {
		if strconv.Itoa(ticket.ID) == id {
			tickets = append(tickets[:i], tickets[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"message": "Ticket deleted"})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
}

// func middlewareLogger() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		c.Next()
// 	}
// }

func middlewareLogger() gin.HandlerFunc {
	// Buka file log untuk menulis atau membuat file baru jika belum ada
	file, err := os.OpenFile("server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Gagal membuka file log:", err)
	}
	// Tetapkan output log ke file
	log.SetOutput(file)

	return func(c *gin.Context) {
		// Tulis log ke file
		log.Printf("[%s] %s %s\n", c.Request.Method, c.Request.URL.Path, c.ClientIP())
		c.Next()
	}
}

func middlewareAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization header"})
			c.Abort()
			return
		}

		token := authHeader         // Mengambil token langsung dari header
		if token != "token_fajar" { // Replace this with your actual token validation logic
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Next()
	}
}
