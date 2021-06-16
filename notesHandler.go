package main

import (
	"context"
	"errors"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
)

type Note struct {
	ID     string `json:"ID"`
	Title  string `json:"Title"`
	Text   string `json:"Text"`
	Author string `json:"-"`
}

type NoteHandler struct {
	Context context.Context
	Client  firestore.Client
}

func NewNoteHandler() (*NoteHandler, error) {
	ctx := context.Background()
	conf := &firebase.Config{ProjectID: "note-taking-app-315314"}
	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		return nil, err
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, err
	}

	return &NoteHandler{
		ctx,
		*client,
	}, nil
}

func NoteWithoutID(note *Note) interface{} {
	return map[string]interface{}{
		"Title":  note.Title,
		"Text":   note.Text,
		"Author": note.Author,
	}
}
func NoteWithoutIDOrAuthor(note *Note) interface{} {
	return map[string]interface{}{
		"Title": note.Title,
		"Text":  note.Text,
	}
}
func (noteHandler *NoteHandler) AddNote(note *Note) (*Note, error) {
	noteCollection := noteHandler.Client.Collection("notes")
	ref := noteCollection.NewDoc()
	note.ID = ref.ID
	_, err := ref.Set(noteHandler.Context, NoteWithoutID(note))
	if err != nil {
		return nil, err
	}

	return note, nil
}

func (noteHandler *NoteHandler) UpdateNote(note *Note) (*Note, error) {
	noteCollection := noteHandler.Client.Collection("notes")
	ref := noteCollection.Doc(note.ID)
	snapshot, err := ref.Get(noteHandler.Context)
	if err != nil {
		return nil, err
	}
	if !snapshot.Exists() {
		return nil, errors.New("file with this ID doesn't exist")
	}

	_, err = ref.Update(noteHandler.Context, []firestore.Update{
		{
			Path:  "Title",
			Value: note.Title,
		},
		{
			Path:  "Text",
			Value: note.Text,
		},
	})
	if err != nil {
		return nil, err
	}

	return note, nil
}

func (noteHandler *NoteHandler) GetNotesByAuthor(author string) ([]Note, error) {
	noteCollection := noteHandler.Client.Collection("notes")

	docIterator := noteCollection.Where("Author", "==", author).Documents(noteHandler.Context)
	docs, err := docIterator.GetAll()
	if err != nil {
		return nil, err
	}

	notes := make([]Note, len(docs))
	for i, docSnapshot := range docs {
		docSnapshot.DataTo(&notes[i])
		notes[i].ID = docSnapshot.Ref.ID
	}
	return notes, nil
}
func (noteHandler *NoteHandler) GetNoteByID(ID string) (*Note, error) {
	noteCollection := noteHandler.Client.Collection("notes")

	docSnapshot, err := noteCollection.Doc(ID).Get(noteHandler.Context)
	if err != nil {
		return nil, err
	}

	var note Note
	err = docSnapshot.DataTo(&note)
	if err != nil {
		return nil, err
	}

	return &note, nil
}

func (noteHandler *NoteHandler) DeleteNote(ID string) error {
	noteCollection := noteHandler.Client.Collection("notes")

	_, err := noteCollection.Doc(ID).Delete(noteHandler.Context)
	return err
}
