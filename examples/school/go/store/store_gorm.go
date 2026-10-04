package store

import (
	"errors"
	"strings"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/school"
)

// StudentRow is the Storage Shape of a Student: it Knows its Table and no one else does.
type StudentRow struct {
	ID       uint   `gorm:"primaryKey"`
	Rut      string `gorm:"uniqueIndex;not null"`
	Name     string `gorm:"not null"`
	Age      int    `gorm:"not null"`
	CourseID uint   `gorm:"not null"`
}

// CourseRow is the Storage Shape of a Course.
type CourseRow struct {
	ID   uint   `gorm:"primaryKey"`
	Code string `gorm:"not null"`
	Name string `gorm:"not null"`
}

// School Keeps Students and Courses in one SQLite database through GORM.
// This Package is the only one that Imports the ORM.
type School struct {
	DB *gorm.DB
}

// OpenSchoolStore Opens the database and Brings its Tables up to Date.
func OpenSchoolStore(path string) (School, error) {
	database, err := gorm.Open(sqlite.Open(path), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		return School{}, err
	}

	if err := database.AutoMigrate(&StudentRow{}, &CourseRow{}); err != nil {
		return School{}, err
	}

	return School{DB: database}, nil
}

// InsertStudentRow Writes a new Student and Returns it Numbered.
// The unique Index Proves the RUT is not Taken.
func (g School) InsertStudentRow(student school.Student) (school.Student, error) {
	row := EncodeStudentRow(student)
	row.ID = 0

	if err := g.DB.Create(&row).Error; err != nil {
		return school.Student{}, TranslateStoreError(err, school.ErrStudentUnknown)
	}

	return DecodeStudentRow(row), nil
}

// SelectStudentRow Finds one Student or Says it is Unknown.
func (g School) SelectStudentRow(id school.StudentID) (school.Student, error) {
	var row StudentRow
	if err := g.DB.First(&row, uint(id)).Error; err != nil {
		return school.Student{}, TranslateStoreError(err, school.ErrStudentUnknown)
	}

	return DecodeStudentRow(row), nil
}

// UpdateStudentRow Rewrites the Student its Id Names, or Says it is Unknown.
func (g School) UpdateStudentRow(student school.Student) (school.Student, error) {
	row := EncodeStudentRow(student)

	result := g.DB.Model(&StudentRow{}).Where("id = ?", row.ID).
		Updates(map[string]any{"rut": row.Rut, "name": row.Name, "age": row.Age, "course_id": row.CourseID})
	if result.Error != nil {
		return school.Student{}, TranslateStoreError(result.Error, school.ErrStudentUnknown)
	}
	if result.RowsAffected == 0 {
		return school.Student{}, school.ErrStudentUnknown
	}

	return student, nil
}

// DeleteStudentRow Removes one Student, or Says it is Unknown.
func (g School) DeleteStudentRow(id school.StudentID) error {
	result := g.DB.Delete(&StudentRow{}, uint(id))
	if result.Error != nil {
		return TranslateStoreError(result.Error, school.ErrStudentUnknown)
	}
	if result.RowsAffected == 0 {
		return school.ErrStudentUnknown
	}

	return nil
}

// SelectStudentPage Reads one Page, Ordered by Name.
func (g School) SelectStudentPage(page school.Page) ([]school.Student, error) {
	var rows []StudentRow
	if err := ApplyPageWindow(g.DB.Order("name").Order("id"), page).Find(&rows).Error; err != nil {
		return nil, TranslateStoreError(err, school.ErrStudentUnknown)
	}

	students := make([]school.Student, 0, len(rows))
	for _, row := range rows {
		students = append(students, DecodeStudentRow(row))
	}

	return students, nil
}

// InsertCourseRow Writes a new Course and Returns it Numbered.
func (g School) InsertCourseRow(course school.Course) (school.Course, error) {
	row := CourseRow{Code: string(course.Code), Name: string(course.Name)}
	if err := g.DB.Create(&row).Error; err != nil {
		return school.Course{}, TranslateStoreError(err, school.ErrCourseUnknown)
	}

	return DecodeCourseRow(row), nil
}

// SelectCourseRow Finds one Course or Says it is Unknown.
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
	if err := ApplyPageWindow(g.DB.Order("code").Order("id"), page).Find(&rows).Error; err != nil {
		return nil, TranslateStoreError(err, school.ErrCourseUnknown)
	}

	courses := make([]school.Course, 0, len(rows))
	for _, row := range rows {
		courses = append(courses, DecodeCourseRow(row))
	}

	return courses, nil
}

// ApplyPageWindow Turns a Page into Limit and Offset; a Size of zero is no Window.
// Reference: https://gorm.io/docs/query.html#Limit-amp-Offset
func ApplyPageWindow(query *gorm.DB, page school.Page) *gorm.DB {
	if page.Size <= 0 {
		return query
	}

	return query.Limit(page.Size).Offset(page.Number * page.Size)
}

// TranslateStoreError Turns a driver Error into a Business Fault when it Means one.
// A driver Error that Means nothing Stays Ours and Reaches the Edge as a 500.
func TranslateStoreError(err error, missing error) error {
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return missing
	case strings.Contains(err.Error(), "UNIQUE constraint failed: student_rows.rut"):
		return school.ErrRutTaken
	default:
		return err
	}
}

// EncodeStudentRow Moves a Business Student into its Storage Shape.
func EncodeStudentRow(student school.Student) StudentRow {
	return StudentRow{
		ID: uint(student.ID), Rut: string(student.Rut), Name: string(student.Name),
		Age: student.Age, CourseID: uint(student.Course),
	}
}

// DecodeStudentRow Moves a Storage row back into the Business Student.
func DecodeStudentRow(row StudentRow) school.Student {
	return school.Student{
		ID: school.StudentID(row.ID), Rut: school.RUT(row.Rut), Name: school.FullName(row.Name),
		Age: row.Age, Course: school.CourseID(row.CourseID),
	}
}

// DecodeCourseRow Moves a Storage row back into the Business Course.
func DecodeCourseRow(row CourseRow) school.Course {
	return school.Course{ID: school.CourseID(row.ID), Code: school.CourseCode(row.Code), Name: school.FullName(row.Name)}
}
