package dto

type PostMessage struct {
	Content         string   `json:"content" binding:"required,min=1,max=8000"`
	MentionEveryone bool     `json:"mention_everyone"`
	Mentions        []string `json:"mentions" binding:"dive,uuid4"`
}
