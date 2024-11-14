package handlers

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/dlrklc/go-note-app/db"
)

type note struct { //todo: delete after db conn
	ID    string `json:"id"`
	Title string `json:"title"`
	Text  string `json:"text"`
}

// todo: Below should get from db
var notes = []note{ //todo: delete after db conn
	{ID: "1", Title: "Hello", Text: "Hello World"},
	{ID: "2", Title: "Test", Text: "Here test for note"},
	{ID: "3", Title: "Drinking coffee", Text: "Should remember to drink coffee"},
}

func GetNotes(c *gin.Context) {
	notes, err := db.GetNotes()
	if err != nil {
		log.Printf("Error getting notes: %v", err)
		return
	}
	c.IndentedJSON(http.StatusOK, notes)
}

func GetNoteByID(c *gin.Context) {
	id := c.Param("id")

	id_i, _ := strconv.Atoi(id)
	note, err := db.GetNote(id_i)
	if err != nil {
		log.Printf("Error getting note: %v", err)
		return
	}
	c.IndentedJSON(http.StatusOK, note)
}

func GetNotesByID(c *gin.Context) {
	idsParam := c.Param("ids")

	ids := strings.Split(idsParam, ",")

	var retrievedNotes = []note{}

	//todo: get post by id from db
	//todo: delete below
	for index, a := range notes {
		for _, b := range ids {
			if a.ID == b {
				retrievedNotes = append(retrievedNotes, notes[index])
				break
			}
		}
	}
	//todo: delete above
	c.IndentedJSON(http.StatusOK, retrievedNotes)
}

func AddNewNote(c *gin.Context) {
	var newNote note

	if err := c.BindJSON(&newNote); err != nil {
		return
	}

	var id, err = db.CreateNote(newNote.Title, newNote.Text)

	if err != nil {
		log.Printf("Error creating note: %v", err)
		return
	}
	c.IndentedJSON(http.StatusCreated, id)
}

func AddNewNotes(c *gin.Context) {
	var newNotes []note

	if err := c.BindJSON(&newNotes); err != nil {
		return
	}
	notes = append(notes, newNotes...) //todo: update in db
	c.IndentedJSON(http.StatusCreated, newNotes)
}

func UpdateNote(c *gin.Context) { //func updateNotes ?
	id := c.Param("id")

	var updatedNote note

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	if err := c.ShouldBindJSON(&updatedNote); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id_i, _ := strconv.Atoi(id)

	err := db.UpdateNote(id_i, updatedNote.Title, updatedNote.Text)

	if err != nil {
		log.Printf("Error updating note: %v", err)
		return
	}

	updatedNote.ID = id
	c.JSON(http.StatusOK, updatedNote)

	/*for i, note := range notes {
		if note.ID == id {
			if updatedNote.Title != "" {
				notes[i].Title = updatedNote.Title
			}
			if updatedNote.Text != "" {
				notes[i].Text = updatedNote.Text
			}
			c.JSON(http.StatusOK, notes[i])
			return
		}
	}

	//todo: update also in db

	c.JSON(http.StatusNotFound, gin.H{"error": "note not found"})*/

}

func DeleteNote(c *gin.Context) {
	id := c.Param("id")

	var deletedNote note

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	for index, a := range notes { //todo: instead get from db
		if a.ID == id {
			deletedNote = a
			notes = append(notes[:index], notes[index+1:]...)
			c.IndentedJSON(http.StatusOK, deletedNote) //todo: also delete from db
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "note not found."})
}

func DeleteNotes(c *gin.Context) {
	idsParam := c.Param("ids")

	ids := strings.Split(idsParam, ",")

	var deletedNotes []note

	if len(ids) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	for _, b := range ids {
		for index, a := range notes {
			if a.ID == b {
				notes = append(notes[:index], notes[index+1:]...)
				deletedNotes = append(deletedNotes, a)
				break
			}
		}
	}

	c.IndentedJSON(http.StatusOK, deletedNotes)
}
