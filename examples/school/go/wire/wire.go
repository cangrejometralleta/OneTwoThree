package wire

// This Package is the Contract, and Imports nothing.
// Reading it Costs no Context.

// StudentBody is what a Client Sends.
// Weak on Purpose: the Wire Cannot be Trusted,
// so nothing here is a Business Type yet.
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
