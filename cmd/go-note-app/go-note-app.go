package main

import (
	"github.com/dlrklc/go-note-app/db"
	"github.com/dlrklc/go-note-app/pkg/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	db.Init()
	defer db.Close()

	router := gin.Default() //create a router

	router.GET("/notes", handlers.GetNotes) //handle albums route
	router.GET("/note/:id", handlers.GetNoteByID)
	router.GET("/notes/:ids", handlers.GetNotesByID)

	router.POST("/note", handlers.AddNewNote)
	router.POST("/notes", handlers.AddNewNotes)

	router.PATCH("/note/:id", handlers.UpdateNote)

	router.DELETE("/note/:id", handlers.DeleteNote)
	router.DELETE("/notes/:ids", handlers.DeleteNotes)

	router.Run("localhost:8080")
}
