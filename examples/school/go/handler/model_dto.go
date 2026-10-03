package handler

import (
	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/school"
	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/wire"
)

func buildStudentRecord(body wire.StudentBody, id school.StudentID) school.Student {
	return school.Student{
		ID: id, Rut: school.RUT(body.Rut), Name: school.FullName(body.Name),
		Age: body.Age, Course: school.CourseID(body.Course),
	}
}

func buildCourseRecord(body wire.CourseBody, id school.CourseID) school.Course {
	return school.Course{ID: id, Code: school.CourseCode(body.Code), Name: school.FullName(body.Name)}
}

func renderStudentView(student school.Student) wire.StudentView {
	return wire.StudentView{
		ID: uint(student.ID), Rut: string(student.Rut), Name: string(student.Name),
		Age: student.Age, Course: uint(student.Course),
	}
}

func renderStudentViews(students []school.Student) []wire.StudentView {
	views := make([]wire.StudentView, 0, len(students))
	for _, student := range students {
		views = append(views, renderStudentView(student))
	}

	return views
}

func renderCourseView(course school.Course) wire.CourseView {
	return wire.CourseView{ID: uint(course.ID), Code: string(course.Code), Name: string(course.Name)}
}

func renderCourseViews(courses []school.Course) []wire.CourseView {
	views := make([]wire.CourseView, 0, len(courses))
	for _, course := range courses {
		views = append(views, renderCourseView(course))
	}

	return views
}
