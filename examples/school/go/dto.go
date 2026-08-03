package main

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

// BuildStudentRecord Promotes untrusted Input into Business Types.
func (b StudentBody) BuildStudentRecord(id StudentID) Student {
	return Student{
		ID:     id,
		Rut:    RUT(b.Rut),
		Name:   FullName(b.Name),
		Age:    b.Age,
		Course: CourseID(b.Course),
	}
}

// BuildCourseRecord Promotes untrusted Course Input.
func (b CourseBody) BuildCourseRecord(id CourseID) Course {
	return Course{ID: id, Code: CourseCode(b.Code), Name: FullName(b.Name)}
}

// RenderStudentView Demotes a Business Student back to the Wire.
func RenderStudentView(s Student) StudentView {
	return StudentView{
		ID:     uint(s.ID),
		Rut:    string(s.Rut),
		Name:   string(s.Name),
		Age:    s.Age,
		Course: uint(s.Course),
	}
}

// RenderCourseView Demotes a Business Course back to the Wire.
func RenderCourseView(c Course) CourseView {
	return CourseView{ID: uint(c.ID), Code: string(c.Code), Name: string(c.Name)}
}

// RenderStudentViews Repeats the Move for a whole Page.
func RenderStudentViews(students []Student) []StudentView {
	views := make([]StudentView, 0, len(students))
	for _, student := range students {
		views = append(views, RenderStudentView(student))
	}
	return views
}

// RenderCourseViews Repeats the Move for a whole Page.
func RenderCourseViews(courses []Course) []CourseView {
	views := make([]CourseView, 0, len(courses))
	for _, course := range courses {
		views = append(views, RenderCourseView(course))
	}
	return views
}
