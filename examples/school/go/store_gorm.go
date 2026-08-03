package main

import (
	"errors"

	"gorm.io/gorm"
)

// StudentRow is the Storage Shape of a Student.
// The Tags Document the Table, so the Schema Explains itself.
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

// GormSchool Fulfils both Stores with one Connection.
// It is the only Type in this Program that Knows GORM Exists.
// Reference: https://gorm.io/docs/
type GormSchool struct {
	DB *gorm.DB
}

// InsertStudentRow Writes a new Student and Returns it Numbered.
func (g GormSchool) InsertStudentRow(s Student) (Student, error) {
	row := EncodeStudentRow(s)

	if err := g.DB.Create(&row).Error; err != nil {
		return Student{}, TranslateStoreError(err, ErrStudentUnknown)
	}

	return DecodeStudentRow(row), nil
}

// SelectStudentRow Finds one Student or Says why not.
func (g GormSchool) SelectStudentRow(id StudentID) (Student, error) {
	var row StudentRow

	if err := g.DB.First(&row, uint(id)).Error; err != nil {
		return Student{}, TranslateStoreError(err, ErrStudentUnknown)
	}

	return DecodeStudentRow(row), nil
}

// UpdateStudentRow Overwrites an existing Student.
func (g GormSchool) UpdateStudentRow(s Student) (Student, error) {
	row := EncodeStudentRow(s)

	result := g.DB.
		Model(&StudentRow{}).
		Where("id = ?", row.ID).
		Updates(row)

	if result.RowsAffected == 0 {
		return Student{}, ErrStudentUnknown
	}

	return s, TranslateStoreError(result.Error, ErrStudentUnknown)
}

// DeleteStudentRow Removes a Student, or Reports an Absence.
func (g GormSchool) DeleteStudentRow(id StudentID) error {
	result := g.DB.Delete(&StudentRow{}, uint(id))

	if result.RowsAffected == 0 {
		return ErrStudentUnknown
	}

	return TranslateStoreError(result.Error, ErrStudentUnknown)
}

// SelectStudentPage Reads one Page, Ordered by Name.
func (g GormSchool) SelectStudentPage(page Page) ([]Student, error) {
	var rows []StudentRow

	query := ApplyPageWindow(g.DB.Order("name"), page)
	if err := query.Find(&rows).Error; err != nil {
		return nil, TranslateStoreError(err, ErrStudentUnknown)
	}

	return DecodeStudentRows(rows), nil
}

// InsertCourseRow Writes a new Course and Returns it Numbered.
func (g GormSchool) InsertCourseRow(c Course) (Course, error) {
	row := CourseRow{ID: uint(c.ID), Code: string(c.Code), Name: string(c.Name)}

	if err := g.DB.Create(&row).Error; err != nil {
		return Course{}, TranslateStoreError(err, ErrCourseUnknown)
	}

	return DecodeCourseRow(row), nil
}

// SelectCourseRow Finds one Course or Says why not.
func (g GormSchool) SelectCourseRow(id CourseID) (Course, error) {
	var row CourseRow

	if err := g.DB.First(&row, uint(id)).Error; err != nil {
		return Course{}, TranslateStoreError(err, ErrCourseUnknown)
	}

	return DecodeCourseRow(row), nil
}

// SelectCoursePage Reads one Page of Courses, Ordered by Code.
func (g GormSchool) SelectCoursePage(page Page) ([]Course, error) {
	var rows []CourseRow

	query := ApplyPageWindow(g.DB.Order("code"), page)
	if err := query.Find(&rows).Error; err != nil {
		return nil, TranslateStoreError(err, ErrCourseUnknown)
	}

	return DecodeCourseRows(rows), nil
}

// ApplyPageWindow Turns a Page into Limit and Offset.
// Reference: https://gorm.io/docs/query.html#Limit-amp-Offset
func ApplyPageWindow(query *gorm.DB, page Page) *gorm.DB {
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

// EncodeStudentRow Turns a Business Student into Storage.
func EncodeStudentRow(s Student) StudentRow {
	return StudentRow{
		ID:       uint(s.ID),
		Rut:      string(s.Rut),
		Name:     string(s.Name),
		Age:      s.Age,
		CourseID: uint(s.Course),
	}
}

// DecodeStudentRow Turns Storage back into Business.
func DecodeStudentRow(r StudentRow) Student {
	return Student{
		ID:     StudentID(r.ID),
		Rut:    RUT(r.Rut),
		Name:   FullName(r.Name),
		Age:    r.Age,
		Course: CourseID(r.CourseID),
	}
}

// DecodeCourseRow Turns a Course Row back into Business.
func DecodeCourseRow(r CourseRow) Course {
	return Course{ID: CourseID(r.ID), Code: CourseCode(r.Code), Name: FullName(r.Name)}
}

// DecodeStudentRows Repeats the Move for a whole Page.
func DecodeStudentRows(rows []StudentRow) []Student {
	students := make([]Student, 0, len(rows))
	for _, row := range rows {
		students = append(students, DecodeStudentRow(row))
	}
	return students
}

// DecodeCourseRows Repeats the Move for a whole Page.
func DecodeCourseRows(rows []CourseRow) []Course {
	courses := make([]Course, 0, len(rows))
	for _, row := range rows {
		courses = append(courses, DecodeCourseRow(row))
	}
	return courses
}
