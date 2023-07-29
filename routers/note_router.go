package routers

import (
	"net/http"

	db "note_app_server/services"

	"github.com/gin-gonic/gin"
)

func GetNotes(c *gin.Context) {
	notes, err := db.GetNotesDB()
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Failed to get notes"})
		return
	}
	c.IndentedJSON(http.StatusOK, notes)
}

func AddNote(c *gin.Context) {
	var newNote db.Note

	if err := c.BindJSON(&newNote); err != nil {
		return
	}

	if err := db.AddNoteDB(newNote); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, newNote)
}

func DeleteNote(c *gin.Context) {
	var deleteNoteID db.NoteDeleteType

	if err := c.BindJSON(&deleteNoteID); err != nil {
		return
	}

	var deleteNote, err = db.GetNoteByIDDB(deleteNoteID.ID)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if deleteNoteID.User != deleteNote.Writer {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "no permission to delete this note"})
		return
	}

	if err := db.DeleteNoteDB(deleteNoteID); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"Deleted note": deleteNote.ID})
}
