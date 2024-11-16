package db

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
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
		return 0, err
	}
	return ID, nil
}

func CreateNotes(notes []Note) ([]int, error) {
	tx, err := DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Prepare insert query
	stmt, err := tx.Prepare(`INSERT INTO notes (title, text) 
		VALUES ($1, $2) RETURNING id`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var ids []int

	for _, note := range notes {
		if note.Title != "" || note.Text != "" {
			var id int
			err := stmt.QueryRow(note.Title, note.Text).Scan(&id)
			if err != nil {
				return nil, err
			}
			ids = append(ids, id)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	// Return the generated ids
	return ids, nil
}

// GetNote retrieves a note by ID
func GetNote(id int) (*Note, error) {
	note := &Note{}
	query := `SELECT id, title, text FROM notes WHERE id = $1`
	err := DB.QueryRow(query, id).Scan(&note.ID, &note.Title, &note.Text)
	if err != nil {
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

func GetNotesByID(ids []int) ([]Note, error) {

	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i := range ids {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = ids[i]
	}

	// Join the placeholders into a comma-separated string
	query := fmt.Sprintf("SELECT id, title, text FROM notes WHERE id IN (%s)", strings.Join(placeholders, ", "))

	// Query the database to retrieve notes
	rows, err := DB.Query(query, args...)
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
func UpdateNote(id, title, text string) (string, error) {
	query := `UPDATE notes SET title = $1, text = $2 WHERE id = $3`
	result, err := DB.Exec(query, title, text, id)
	if err != nil {
		return "", err
	}
	message := createResultMessage(&result, "updated")
	return message, err
}

func UpdateNotes(notes []Note) ([]string, error) {
	// Begin a transaction
	tx, err := DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Prepare the update statement
	stmt, err := tx.Prepare(`UPDATE notes SET title = $1, text = $2 WHERE id = $3`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Update each note in the slice
	var message []string
	for _, note := range notes {
		if (note.Title != "" || note.Text != "") && note.ID != "" {
			result, err := stmt.Exec(note.Title, note.Text, note.ID)
			if err != nil {
				return nil, err
			}
			message = append(message, createResultMessage(&result, "updated"))
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return message, nil
}

// DeleteNote deletes a note by ID
func DeleteNote(id int) (string, error) {
	query := `DELETE FROM notes WHERE id = $1`
	result, err := DB.Exec(query, id)

	message := createResultMessage(&result, "deleted")

	return message, err
}

// DeleteNote deletes a note by ID
func DeleteNotes(ids []int) (string, error) {
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i := range ids {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = ids[i]
	}

	query := fmt.Sprintf("DELETE FROM notes WHERE id IN (%s)", strings.Join(placeholders, ", "))
	result, err := DB.Exec(query, args...)

	message := createResultMessage(&result, "deleted")

	return message, err
}

func createResultMessage(result *sql.Result, messageType string) string {
	// Get the number of rows affected
	rowsAffected, err := (*result).RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	var message string
	if rowsAffected > 0 {
		message = fmt.Sprintf("Successfully %s %d row(s)", messageType, rowsAffected)
	} else {
		message = fmt.Sprintf("No rows were %s (note not found)", messageType)
	}
	return message
}
