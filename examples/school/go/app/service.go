package app

import "github.com/cangrejometralleta/OneTwoThree/examples/school/go/school"

// SchoolService Composes the Stores behind each School Use Case.
type SchoolService struct {
	Students StudentStore
	Courses  CourseStore
}

// ListStudents Reads one Page of the Roll.
func (s SchoolService) ListStudents(page school.Page) ([]school.Student, error) {
	return s.Students.SelectStudentPage(page)
}

// ReadStudent Finds one Enrolment.
func (s SchoolService) ReadStudent(id school.StudentID) (school.Student, error) {
	return s.Students.SelectStudentRow(id)
}

// EnrollStudent Validates and Saves one Enrolment.
func (s SchoolService) EnrollStudent(student school.Student) (school.Student, error) {
	if err := s.checkStudent(student); err != nil {
		return school.Student{}, err
	}

	return s.Students.InsertStudentRow(student)
}

// SaveStudent Validates and Rewrites one Enrolment.
func (s SchoolService) SaveStudent(student school.Student) (school.Student, error) {
	if err := s.checkStudent(student); err != nil {
		return school.Student{}, err
	}

	return s.Students.UpdateStudentRow(student)
}

// DeleteStudent Removes one Enrolment.
func (s SchoolService) DeleteStudent(id school.StudentID) error {
	return s.Students.DeleteStudentRow(id)
}

// ListCourses Reads one Page of the Catalogue.
func (s SchoolService) ListCourses(page school.Page) ([]school.Course, error) {
	return s.Courses.SelectCoursePage(page)
}

// ReadCourse Finds one Course.
func (s SchoolService) ReadCourse(id school.CourseID) (school.Course, error) {
	return s.Courses.SelectCourseRow(id)
}

// CreateCourse Stores one Course.
func (s SchoolService) CreateCourse(course school.Course) (school.Course, error) {
	return s.Courses.InsertCourseRow(course)
}

func (s SchoolService) checkStudent(student school.Student) error {
	if err := student.CheckStudentRecord(); err != nil {
		return err
	}

	_, err := s.Courses.SelectCourseRow(student.Course)
	return err
}
