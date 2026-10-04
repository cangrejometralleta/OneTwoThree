package wire

// This Package is the Contract, and Imports nothing.
// Open one File and you Know the whole API.

// StudentBody is what a Client Sends. Weak on Purpose: the Wire cannot be
// Trusted, so nothing here is a Business Type yet.
type StudentBody struct {
	Rut    string `json:"rut"`
	Name   string `json:"name"`
	Age    int    `json:"age"`
	Course uint   `json:"courseId"`
}

// StudentView is what a Client Gets back.
type StudentView struct {
	ID     uint   `json:"id"`
	Rut    string `json:"rut"`
	Name   string `json:"name"`
	Age    int    `json:"age"`
	Course uint   `json:"courseId"`
}

// CourseBody is what a Client Sends about a Course.
type CourseBody struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// CourseView is what a Client Gets back about a Course.
type CourseView struct {
	ID   uint   `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

// TokenView is what the Door Hands out for a Token.
type TokenView struct {
	Token string `json:"token"`
}

// FailureView is the one Shape every Failure Answers with.
type FailureView struct {
	Error string `json:"error"`
}
