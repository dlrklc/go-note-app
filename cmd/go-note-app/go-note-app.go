package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default() //create a router

	/*router.GET("/notes", getNotes) //handle albums route
	router.GET("/notes/:id", getNoteByID)
	router.POST("/notes", addNewNote)
	router.PATCH("/notes/:id", updateNote)*/

	router.Run("localhost:8080")
}
