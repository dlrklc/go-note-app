package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type note struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Text  string `json:"text"`
}

// todo: Below should get from db
var notes = []note{
	{ID: "1", Title: "Hello", Text: "Hello World"},
	{ID: "2", Title: "Test", Text: "Here test for note"},
	{ID: "3", Title: "Drinking coffee", Text: "Should remember to drink coffee"},
}

func getNotes(c *gin.Context) {
	//todo: get posts from db
	c.IndentedJSON(http.StatusOK, notes)
}

func getNoteByID(c *gin.Context) {
	id := c.Param("id")
	//todo: get post by id from db
	//todo: delete below
	for _, a := range notes {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	//todo: delete above
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "note not found"})
}

func addNewNote(c *gin.Context) {
	var newNote note

	if err := c.BindJSON(&newNote); err != nil {
		return
	}
	notes = append(notes, newNote) //todo: update in db
	c.IndentedJSON(http.StatusCreated, newNote)
}

func updateNote(c *gin.Context) {
	id, ok := c.GetQuery("id")

	var found note

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id query parameter."})
		return
	}

	for _, a := range notes { //todo: instead get from db
		if a.ID == id {
			found = a
		}
	}

	if found == (note{}) {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "note not found."})
		return
	}

	title := c.Param("title")
	text := c.Param("text")

	found.Title = title
	found.Text = text

	//todo: update found in db

	c.IndentedJSON(http.StatusOK, found)

}

func deleteNote(c *gin.Context) {
	id, ok := c.GetQuery("id")

	var found note

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id query parameter."})
		return
	}

	for _, a := range notes { //todo: instead get from db
		if a.ID == id {
			found = a
		}
	}

	if found == (note{}) {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "note not found."})
		return
	}

	//delete from db

	c.IndentedJSON(http.StatusOK, found)
}

/*func deleteNotes
func updateNotesnote
func addNotes ?
func getNotesByID*/
