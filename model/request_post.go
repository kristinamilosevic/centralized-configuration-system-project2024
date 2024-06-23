package model

// swagger:model RequestPost
type RequestPost struct {
	// Id of the post
	// in: string
	Id string `json:"id"`

	// Title of the post
	// in: string
	Title string `json:"title"`

	// Text content of the post
	// in: string
	Text string `json:"text"`

	// List of tags of the post
	// in: []string
	Tags []string `json:"tags"`
}
