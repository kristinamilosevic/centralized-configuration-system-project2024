package model

// swagger:parameters deletePost
type DeleteRequest struct {
	// Post ID
	// in: path
	Id string `json:"id"`
}

// swagger:parameters getPostById
type GetRequest struct {
	// Post ID
	// in: path
	Id string `json:"id"`
}

// swagger:parameters post createPost
type RequestPostBody struct {
	// - name: body
	//  in: body
	//  description: name and status
	//  schema:
	//  type: object
	//     "$ref": "#/definitions/RequestPost"
	//  required: true
	Body RequestPost `json:"body"`
}
