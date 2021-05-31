package models

type ChangeConsumer interface {
	ModelChanged(id string, senderId *string)
}

var (
	Consumer ChangeConsumer
)

type DocumentWorker struct{}

func (dw DocumentWorker) FetchDocument(resourceId string) (string, error) {
	note, err := FetchNote("", resourceId)
	if note == nil || err != nil {
		return "", err
	}

	return note.Content, err
}

func (dw DocumentWorker) SaveDocument(resourceId string, content string) {
	SaveNoteContent(resourceId, content)
}
