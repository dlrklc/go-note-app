package db

import (
	"database/sql"
	"fmt"
	"log"
)

type Note struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Text  string `json:"text"`
}

// CreateNote inserts a new note into the database
func CreateNote(title, text string) (int, error) {
	var ID int
	query := `INSERT INTO notes (title, text) VALUES ($1, $2) RETURNING id`
	err := DB.QueryRow(query, title, text).Scan(&ID)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return 0, err
	}
	return ID, nil
}

// GetNote retrieves a note by ID
func GetNote(id int) (*Note, error) {
	note := &Note{}
	query := `SELECT id, title, text FROM notes WHERE id = $1`
	err := DB.QueryRow(query, id).Scan(&note.ID, &note.Title, &note.Text)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("note not found")
		}
		return nil, err
	}
	return note, nil
}

func GetNotes() ([]Note, error) {
	query := `SELECT id, title, text FROM notes`

	// Query the database to retrieve notes
	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []Note
	for rows.Next() {
		var note Note
		// Scan the row into the Note struct
		err := rows.Scan(&note.ID, &note.Title, &note.Text)
		if err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}

	// Check for any error from iterating over rows
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return notes, nil
}

// UpdateNote updates note information
func UpdateNote(id int, title, text string) error {
	query := `UPDATE notes SET title = $1, text = $2 WHERE id = $3`
	_, err := DB.Exec(query, title, text, id)
	return err
}

// DeleteNote deletes a note by ID
func DeleteNote(id int) error {
	query := `DELETE FROM notes WHERE id = $1`
	_, err := DB.Exec(query, id)
	return err
}
