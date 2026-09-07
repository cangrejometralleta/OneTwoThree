package school

import "github.com/cangrejometralleta/OneTwoThree/examples/school/go/wire"

// The Wire Package Holds Shapes; this File Holds the Crossing.
// Promotion Lives with the Business Types it Promotes into.

// BuildStudentRecord Promotes untrusted Input into Business Types.
func BuildStudentRecord(body wire.StudentBody, id StudentID) Student {
	return Student{
		ID:     id,
		Rut:    RUT(body.Rut),
		Name:   FullName(body.Name),
		Age:    body.Age,
		Course: CourseID(body.Course),
	}
}

// BuildCourseRecord Promotes untrusted Course Input.
func BuildCourseRecord(body wire.CourseBody, id CourseID) Course {
	return Course{ID: id, Code: CourseCode(body.Code), Name: FullName(body.Name)}
}

// RenderStudentView Demotes a Business Student back to the Wire.
func RenderStudentView(s Student) wire.StudentView {
	return wire.StudentView{
		ID:     uint(s.ID),
		Rut:    string(s.Rut),
		Name:   string(s.Name),
		Age:    s.Age,
		Course: uint(s.Course),
	}
}

// RenderCourseView Demotes a Business Course back to the Wire.
func RenderCourseView(c Course) wire.CourseView {
	return wire.CourseView{ID: uint(c.ID), Code: string(c.Code), Name: string(c.Name)}
}

// RenderStudentViews Repeats the Move for a whole Page.
func RenderStudentViews(students []Student) []wire.StudentView {
	views := make([]wire.StudentView, 0, len(students))
	for _, student := range students {
		views = append(views, RenderStudentView(student))
	}

	return views
}

// RenderCourseViews Repeats the Move for a whole Page.
func RenderCourseViews(courses []Course) []wire.CourseView {
	views := make([]wire.CourseView, 0, len(courses))
	for _, course := range courses {
		views = append(views, RenderCourseView(course))
	}

	return views
}
