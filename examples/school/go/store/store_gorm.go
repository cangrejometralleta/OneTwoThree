package store

import (
	"errors"
	"strings"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/school"
)

// StudentRow is the Storage Shape of a Student.
// The Tags Document the Table.
// Reference: https://gorm.io/docs/models.html
type StudentRow struct {
	ID       uint   `gorm:"primaryKey"`
	Rut      string `gorm:"size:16;uniqueIndex;not null"`
	Name     string `gorm:"size:120;not null"`
	Age      int    `gorm:"not null"`
	CourseID uint   `gorm:"index;not null"`
}

// CourseRow is the Storage Shape of a Course.
type CourseRow struct {
	ID   uint   `gorm:"primaryKey"`
	Code string `gorm:"size:16;uniqueIndex;not null"`
	Name string `gorm:"size:120;not null"`
}

// School Fulfils both Stores with one Connection.
// It is the only Type in this Program that Knows GORM Exists.
// Reference: https://gorm.io/docs/
type School struct {
	DB *gorm.DB
}

// OpenSchoolStore Opens SQLite, Shapes both Tables and Hands back the Store.
// The Caller Decides what a Failure Means.
// Reference: https://gorm.io/docs/connecting_to_the_database.html
func OpenSchoolStore(path string) (School, error) {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return School{}, err
	}

	if err := db.AutoMigrate(&StudentRow{}, &CourseRow{}); err != nil {
		return School{}, err
	}

	return School{DB: db}, nil
}

// InsertStudentRow Writes a new Student and Returns it Numbered.
func (g School) InsertStudentRow(s school.Student) (school.Student, error) {
	row := EncodeStudentRow(s)

	if err := g.DB.Create(&row).Error; err != nil {
		return school.Student{}, TranslateInsertError(err)
	}

	return DecodeStudentRow(row), nil
}

// SelectStudentRow Finds one Student or Says why not.
func (g School) SelectStudentRow(id school.StudentID) (school.Student, error) {
	var row StudentRow

	if err := g.DB.First(&row, uint(id)).Error; err != nil {
		return school.Student{}, TranslateStoreError(err, school.ErrStudentUnknown)
	}

	return DecodeStudentRow(row), nil
}

// UpdateStudentRow Overwrites an existing Student.
func (g School) UpdateStudentRow(s school.Student) (school.Student, error) {
	row := EncodeStudentRow(s)

	result := g.DB.
		Model(&StudentRow{}).
		Where("id = ?", row.ID).
		Updates(row)

	if result.RowsAffected == 0 {
		return school.Student{}, school.ErrStudentUnknown
	}

	return s, TranslateStoreError(result.Error, school.ErrStudentUnknown)
}

// DeleteStudentRow Removes a Student, or Reports an Absence.
func (g School) DeleteStudentRow(id school.StudentID) error {
	result := g.DB.Delete(&StudentRow{}, uint(id))

	if result.RowsAffected == 0 {
		return school.ErrStudentUnknown
	}

	return TranslateStoreError(result.Error, school.ErrStudentUnknown)
}

// SelectStudentPage Reads one Page, Ordered by Name.
func (g School) SelectStudentPage(page school.Page) ([]school.Student, error) {
	var rows []StudentRow

	query := ApplyPageWindow(g.DB.Order("name"), page)
	if err := query.Find(&rows).Error; err != nil {
		return nil, TranslateStoreError(err, school.ErrStudentUnknown)
	}

	return DecodeStudentRows(rows), nil
}

// InsertCourseRow Writes a new Course and Returns it Numbered.
func (g School) InsertCourseRow(c school.Course) (school.Course, error) {
	row := CourseRow{ID: uint(c.ID), Code: string(c.Code), Name: string(c.Name)}

	if err := g.DB.Create(&row).Error; err != nil {
		return school.Course{}, TranslateStoreError(err, school.ErrCourseUnknown)
	}

	return DecodeCourseRow(row), nil
}

// SelectCourseRow Finds one Course or Says why not.
func (g School) SelectCourseRow(id school.CourseID) (school.Course, error) {
	var row CourseRow

	if err := g.DB.First(&row, uint(id)).Error; err != nil {
		return school.Course{}, TranslateStoreError(err, school.ErrCourseUnknown)
	}

	return DecodeCourseRow(row), nil
}

// SelectCoursePage Reads one Page of Courses, Ordered by Code.
func (g School) SelectCoursePage(page school.Page) ([]school.Course, error) {
	var rows []CourseRow

	query := ApplyPageWindow(g.DB.Order("code"), page)
	if err := query.Find(&rows).Error; err != nil {
		return nil, TranslateStoreError(err, school.ErrCourseUnknown)
	}

	return DecodeCourseRows(rows), nil
}

// ApplyPageWindow Turns a Page into Limit and Offset.
// Reference: https://gorm.io/docs/query.html#Limit-amp-Offset
func ApplyPageWindow(query *gorm.DB, page school.Page) *gorm.DB {
	if page.Size <= 0 {
		return query
	}

	return query.Offset(page.Number * page.Size).Limit(page.Size)
}

// TranslateStoreError Turns a Driver Failure into a Business one.
func TranslateStoreError(err error, absent error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return absent
	}

	return err
}

// TranslateInsertError Names the one Way a Write can Collide.
// The Message Text is the only Signal SQLite Gives, in either Driver.
func TranslateInsertError(err error) error {
	if err != nil && strings.Contains(err.Error(), "UNIQUE constraint failed") {
		return school.ErrRutTaken
	}

	return err
}

// EncodeStudentRow Turns a Business Student into Storage.
func EncodeStudentRow(s school.Student) StudentRow {
	return StudentRow{
		ID:       uint(s.ID),
		Rut:      string(s.Rut),
		Name:     string(s.Name),
		Age:      s.Age,
		CourseID: uint(s.Course),
	}
}

// DecodeStudentRow Turns Storage back into Business.
func DecodeStudentRow(r StudentRow) school.Student {
	return school.Student{
		ID:     school.StudentID(r.ID),
		Rut:    school.RUT(r.Rut),
		Name:   school.FullName(r.Name),
		Age:    r.Age,
		Course: school.CourseID(r.CourseID),
	}
}

// DecodeCourseRow Turns a Course Row back into Business.
func DecodeCourseRow(r CourseRow) school.Course {
	return school.Course{ID: school.CourseID(r.ID), Code: school.CourseCode(r.Code), Name: school.FullName(r.Name)}
}

// DecodeStudentRows Repeats the Move for a whole Page.
func DecodeStudentRows(rows []StudentRow) []school.Student {
	students := make([]school.Student, 0, len(rows))
	for _, row := range rows {
		students = append(students, DecodeStudentRow(row))
	}
	return students
}

// DecodeCourseRows Repeats the Move for a whole Page.
func DecodeCourseRows(rows []CourseRow) []school.Course {
	courses := make([]school.Course, 0, len(rows))
	for _, row := range rows {
		courses = append(courses, DecodeCourseRow(row))
	}
	return courses
}
