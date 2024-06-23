package model

// swagger:response ResponsePost
type ResponsePost struct {
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

// swagger:response ErrorResponse
type ErrorResponse struct {
	// Error status code
	// in: int64
	Status int64 `json:"status"`
	// Message of the error
	// in: string
	Message string `json:"message"`
}

// swagger:response NoContentResponse
type NoContentResponse struct{}
