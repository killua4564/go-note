package note

type Account struct {
	Username string `json:"username"`
}

type Note struct {
	SID     string `json:"sid"`
	Topic   string `json:"topic"`
	Content string `json:"content,omitempty"`
}

type NoteList struct {
	Viewer []Note `json:"viewer"`
	Owner  []Note `json:"owner"`
}

type AccountNote struct {
	AccountID int64 `json:"account_id"`
	NoteID    int64 `json:"note_id"`
	IsOwner   int64 `json:"is_owner"`
}

type validateNote struct {
	Topic   string `json:"topic" binding:"required"`
	Content string `json:"content"`
}

type validateNoteUri struct {
	SID string `uri:"sid" binding:"required,uuid"`
}

type validateAccountNote struct {
	Username string `json:"username" binding:"required"`
	IsOwner  *int64 `json:"is_owner" binding:"required"`
}
