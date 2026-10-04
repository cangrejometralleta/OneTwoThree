package app

import (
	"errors"
	"sort"
	"testing"

	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/school"
)

// shelf Stands in for the Store: Course 1 is Open, and it Keeps what it is Given.
type shelf struct {
	students map[school.StudentID]school.Student
	courses  map[school.CourseID]school.Course
	page     school.Page
}

func buildShelf() *shelf {
	return &shelf{
		students: map[school.StudentID]school.Student{},
		courses:  map[school.CourseID]school.Course{1: {ID: 1, Code: "MAT-101", Name: "Algebra"}},
	}
}

func (s *shelf) InsertStudentRow(student school.Student) (school.Student, error) {
	student.ID = school.StudentID(len(s.students) + 1)
	s.students[student.ID] = student

	return student, nil
}

func (s *shelf) SelectStudentRow(id school.StudentID) (school.Student, error) {
	student, found := s.students[id]
	if !found {
		return school.Student{}, school.ErrStudentUnknown
	}

	return student, nil
}

func (s *shelf) UpdateStudentRow(student school.Student) (school.Student, error) {
	if _, found := s.students[student.ID]; !found {
		return school.Student{}, school.ErrStudentUnknown
	}
	s.students[student.ID] = student

	return student, nil
}

func (s *shelf) DeleteStudentRow(id school.StudentID) error {
	if _, found := s.students[id]; !found {
		return school.ErrStudentUnknown
	}
	delete(s.students, id)

	return nil
}

func (s *shelf) SelectStudentPage(page school.Page) ([]school.Student, error) {
	s.page = page
	roll := make([]school.Student, 0, len(s.students))
	for _, student := range s.students {
		roll = append(roll, student)
	}
	sort.Slice(roll, func(a, b int) bool { return roll[a].ID < roll[b].ID })

	return roll, nil
}

func (s *shelf) SelectCourseRow(id school.CourseID) (school.Course, error) {
	course, found := s.courses[id]
	if !found {
		return school.Course{}, school.ErrCourseUnknown
	}

	return course, nil
}

func (s *shelf) InsertCourseRow(course school.Course) (school.Course, error) {
	course.ID = school.CourseID(len(s.courses) + 1)
	s.courses[course.ID] = course

	return course, nil
}

func (s *shelf) SelectCoursePage(page school.Page) ([]school.Course, error) {
	s.page = page
	catalogue := make([]school.Course, 0, len(s.courses))
	for _, course := range s.courses {
		catalogue = append(catalogue, course)
	}
	sort.Slice(catalogue, func(a, b int) bool { return catalogue[a].ID < catalogue[b].ID })

	return catalogue, nil
}

// office Answers both Questions, the Way one Adapter may Fill two Ports.
type office struct {
	registered bool
	down       error
	noticeErr  error
	asked      []school.RUT
	announced  []school.Student
}

func (o *office) ConfirmRut(rut school.RUT) (bool, error) {
	o.asked = append(o.asked, rut)

	return o.registered, o.down
}

func (o *office) AnnounceEnrollment(student school.Student) error {
	o.announced = append(o.announced, student)

	return o.noticeErr
}

func buildAdmission(front *office) (SchoolService, *shelf) {
	books := buildShelf()

	return SchoolService{Students: books, Courses: books, Registry: front, Notices: front}, books
}

func readyStudent() school.Student {
	return school.Student{Rut: "11111111-1", Name: "Ana", Age: 20, Course: 1}
}

func TestEnrollStudentAnnouncesOnceTheStoreKeeps(t *testing.T) {
	front := &office{registered: true}
	service, books := buildAdmission(front)

	kept, err := service.EnrollStudent(readyStudent())

	if err != nil || len(books.students) != 1 || len(front.announced) != 1 || front.announced[0].ID != kept.ID {
		t.Fatalf("kept=%v err=%v stored=%d announced=%d", kept, err, len(books.students), len(front.announced))
	}
}

func TestEnrollStudentRefusesWhatTheRegistryDenies(t *testing.T) {
	front := &office{registered: false}
	service, books := buildAdmission(front)

	_, err := service.EnrollStudent(readyStudent())

	if !errors.Is(err, school.ErrRutUnregistered) || len(books.students) != 0 || len(front.announced) != 0 {
		t.Fatalf("err=%v stored=%d announced=%d", err, len(books.students), len(front.announced))
	}
}

func TestEnrollStudentNeverTreatsAnUnknownAsAYes(t *testing.T) {
	front := &office{registered: true, down: errors.New("socket closed")}
	service, books := buildAdmission(front)

	_, err := service.EnrollStudent(readyStudent())

	if !errors.Is(err, school.ErrRegistryUnavailable) || len(books.students) != 0 {
		t.Fatalf("err=%v stored=%d", err, len(books.students))
	}

	service.Registry = nil
	if _, err := service.EnrollStudent(readyStudent()); !errors.Is(err, school.ErrRegistryUnavailable) {
		t.Fatalf("a missing Registry must not Pass, err=%v", err)
	}
}

func TestEnrollStudentKeepsAFailedNoticeInsideTheAnswer(t *testing.T) {
	front := &office{registered: true, noticeErr: errors.New("mail is down")}
	service, books := buildAdmission(front)

	_, err := service.EnrollStudent(readyStudent())

	if err != nil || len(books.students) != 1 {
		t.Fatalf("a Notice that did not Leave must not Fail the Enrolment, err=%v stored=%d", err, len(books.students))
	}
}

func TestTheRegistryIsAskedAfterTheFormAndTheCourse(t *testing.T) {
	front := &office{registered: true}
	service, _ := buildAdmission(front)

	invalid := readyStudent()
	invalid.Name = ""
	unknownCourse := readyStudent()
	unknownCourse.Course = 9

	_, formErr := service.EnrollStudent(invalid)
	_, courseErr := service.EnrollStudent(unknownCourse)

	if !errors.Is(formErr, school.ErrNameIsEmpty) || !errors.Is(courseErr, school.ErrCourseUnknown) || len(front.asked) != 0 {
		t.Fatalf("form=%v course=%v asked=%d", formErr, courseErr, len(front.asked))
	}
}

func TestSaveStudentAsksTheRegistryAndAnnouncesNothing(t *testing.T) {
	front := &office{registered: true}
	service, _ := buildAdmission(front)
	kept, _ := service.EnrollStudent(readyStudent())
	announcedByEnrolment := len(front.announced)

	front.registered = false
	if _, err := service.SaveStudent(kept); !errors.Is(err, school.ErrRutUnregistered) {
		t.Fatalf("a rewrite must ask the Registry too, err=%v", err)
	}

	front.registered = true
	if _, err := service.SaveStudent(kept); err != nil || len(front.announced) != announcedByEnrolment {
		t.Fatalf("a rewrite announces nothing, err=%v announced=%d", err, len(front.announced))
	}
}

func TestSaveStudentRefusesAStudentNobodyHolds(t *testing.T) {
	service, _ := buildAdmission(&office{registered: true})
	ghost := readyStudent()
	ghost.ID = 42

	if _, err := service.SaveStudent(ghost); !errors.Is(err, school.ErrStudentUnknown) {
		t.Fatalf("err=%v", err)
	}
}

// The plain Use Cases, each Told through the Store Ports and nothing else.

func TestReadStudentFindsWhatEnrolmentKept(t *testing.T) {
	service, _ := buildAdmission(&office{registered: true})
	kept, _ := service.EnrollStudent(readyStudent())

	found, err := service.ReadStudent(kept.ID)

	if err != nil || found != kept {
		t.Fatalf("found=%v kept=%v err=%v", found, kept, err)
	}
}

func TestReadStudentRefusesWhoNeverEnrolled(t *testing.T) {
	service, _ := buildAdmission(&office{registered: true})

	if _, err := service.ReadStudent(42); !errors.Is(err, school.ErrStudentUnknown) {
		t.Fatalf("err=%v", err)
	}
}

func TestListStudentsReadsTheRequestedPageOfTheRoll(t *testing.T) {
	service, books := buildAdmission(&office{registered: true})
	service.EnrollStudent(readyStudent())

	roll, err := service.ListStudents(school.Page{Number: 2, Size: 5})

	if err != nil || len(roll) != 1 || books.page != (school.Page{Number: 2, Size: 5}) {
		t.Fatalf("roll=%v page=%v err=%v", roll, books.page, err)
	}
}

func TestDeleteStudentEndsAnEnrolmentOnce(t *testing.T) {
	service, _ := buildAdmission(&office{registered: true})
	kept, _ := service.EnrollStudent(readyStudent())

	if err := service.DeleteStudent(kept.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := service.ReadStudent(kept.ID); !errors.Is(err, school.ErrStudentUnknown) {
		t.Fatalf("a Student who Left must be Unknown, err=%v", err)
	}
	if err := service.DeleteStudent(kept.ID); !errors.Is(err, school.ErrStudentUnknown) {
		t.Fatalf("dropping twice must be Unknown, err=%v", err)
	}
}

func TestCreateCourseOpensACourseTheCatalogueLists(t *testing.T) {
	service, books := buildAdmission(&office{registered: true})

	opened, err := service.CreateCourse(school.Course{Code: "FIS-201", Name: "Physics"})
	if err != nil || opened.ID == 0 {
		t.Fatalf("opened=%v err=%v", opened, err)
	}

	found, _ := service.ReadCourse(opened.ID)
	catalogue, _ := service.ListCourses(school.Page{Number: 0, Size: 10})

	if found != opened || len(catalogue) != 2 || books.page.Size != 10 {
		t.Fatalf("found=%v catalogue=%v page=%v", found, catalogue, books.page)
	}
}

func TestCreateCourseRefusesACourseWithoutCodeOrName(t *testing.T) {
	service, books := buildAdmission(&office{registered: true})

	_, noCode := service.CreateCourse(school.Course{Name: "Physics"})
	_, noName := service.CreateCourse(school.Course{Code: "FIS-201"})

	if !errors.Is(noCode, school.ErrCodeIsEmpty) || !errors.Is(noName, school.ErrNameIsEmpty) || len(books.courses) != 1 {
		t.Fatalf("noCode=%v noName=%v courses=%d", noCode, noName, len(books.courses))
	}
}

func TestReadCourseRefusesACourseThatNeverOpened(t *testing.T) {
	service, _ := buildAdmission(&office{registered: true})

	if _, err := service.ReadCourse(99); !errors.Is(err, school.ErrCourseUnknown) {
		t.Fatalf("err=%v", err)
	}
}
