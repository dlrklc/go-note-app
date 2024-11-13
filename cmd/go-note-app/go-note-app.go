package main

import (
	"github.com/dlrklc/go-note-app/pkg/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default() //create a router

	router.GET("/notes", handlers.GetNotes) //handle albums route
	router.GET("/notes/:id", handlers.GetNoteByID)
	router.POST("/notes", handlers.AddNewNote)
	router.PATCH("/notes/:id", handlers.UpdateNote)

	router.Run("localhost:8080")
}
