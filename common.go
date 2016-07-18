package box

type User struct {
	Type  string `json:"type"`
	ID    string `json:"id"`
	Name  string `json:"name"`
	Login string `json:"login"`
}

type ItemParent struct {
	Type       string `json:"type"`
	ID         string `json:"id"`
	SequenceID string `json:"sequence_id"`
	Etag       string `json:"etag"`
	Name       string `json:"name"`
}
