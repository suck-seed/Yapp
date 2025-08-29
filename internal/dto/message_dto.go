package dto

type PostMessage struct {
	Content         string   `json:"content" validate:"required,min=1,max=8000"`
	MentionEveryone bool     `json:"mention_everyone"`
	Mentions        []string `json:"mentions" validate:"dive,uuid4"`
}
