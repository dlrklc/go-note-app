package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/dlrklc/go-note-app/db"
)

func GetNotes(c *gin.Context) {
	notes, err := db.GetNotes()
	if err != nil {
		log.Printf("[GetNotes] Error getting notes: %v", err)
		if err == sql.ErrNoRows {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "note not found"})
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.IndentedJSON(http.StatusOK, notes)
}

func GetNoteByID(c *gin.Context) {
	id := c.Param("id")

	id_i, _ := strconv.Atoi(id)
	note, err := db.GetNote(id_i)
	if err != nil {
		log.Printf("[GetNoteByID] Error getting note: %v", err)
		if err == sql.ErrNoRows {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "note not found"})
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.IndentedJSON(http.StatusOK, note)
}

func GetNotesByID(c *gin.Context) {
	idsParam := c.Param("ids")

	ids := strings.Split(idsParam, ",")

	var ids_i []int

	for _, id := range ids {
		id_i, err := strconv.Atoi(id)
		if err == nil {
			ids_i = append(ids_i, id_i)
		}
	}

	if len(ids_i) == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "note not found"})
		return
	}

	notes, err := db.GetNotesByID(ids_i)
	if err != nil {
		log.Printf("[GetNotesByID] Error getting notes: %v", err)
		c.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}

	if len(notes) == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "note not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, notes)
}

func AddNewNote(c *gin.Context) {
	var newNote db.Note

	if err := c.BindJSON(&newNote); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}

	if newNote.Title == "" && newNote.Text == "" {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "note must have title or text"})
		return
	}

	var id, err = db.CreateNote(newNote.Title, newNote.Text)

	if err != nil {
		log.Printf("[AddNewNote] Error creating note: %v", err)
		c.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.IndentedJSON(http.StatusCreated, gin.H{"id": id})
}

func AddNewNotes(c *gin.Context) {
	var newNotes []db.Note

	if err := c.BindJSON(&newNotes); err != nil {
		return
	}

	var ids []int
	ids, err := db.CreateNotes(newNotes)
	if err != nil {
		log.Printf("[AddNewNotes] Error creating notes: %v", err)
		c.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}

	if len(ids) == 0 {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "note content cannot be empty"})
		return
	}

	c.IndentedJSON(http.StatusCreated, gin.H{"ids": ids})
}

func UpdateNote(c *gin.Context) {

	var updatedNote db.Note

	if err := c.BindJSON(&updatedNote); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if updatedNote.ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	if updatedNote.Title == "" && updatedNote.Text == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "note must have title or text"})
		return
	}

	message, err := db.UpdateNote(updatedNote.ID, updatedNote.Title, updatedNote.Text)

	if err != nil {
		log.Printf("[UpdateNote] Error updating note: %v", err)
		c.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, message)

}

func UpdateNotes(c *gin.Context) {
	var updatedNotes []db.Note

	if err := c.BindJSON(&updatedNotes); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}

	message, err := db.UpdateNotes(updatedNotes)
	if err != nil {
		log.Printf("[UpdateNotes] Error updating notes: %v", err)
		c.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}

	if len(message) == 0 {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "note must have id and title or text"})
		return
	}

	c.IndentedJSON(http.StatusCreated, message)
}

func DeleteNote(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	id_i, _ := strconv.Atoi(id)

	message, err := db.DeleteNote(id_i)

	if err != nil {
		log.Printf("[DeleteNote] Error deleting note: %v", err)
		c.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.IndentedJSON(http.StatusOK, message)
}

func DeleteNotes(c *gin.Context) {
	idsParam := c.Param("ids")

	ids := strings.Split(idsParam, ",")

	if len(ids) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	var ids_i []int

	for _, id := range ids {
		id_i, _ := strconv.Atoi(id)
		ids_i = append(ids_i, id_i)
	}

	message, err := db.DeleteNotes(ids_i)

	if err != nil {
		log.Printf("[DeleteNotes] Error deleting notes: %v", err)
		c.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.IndentedJSON(http.StatusOK, message)
}
