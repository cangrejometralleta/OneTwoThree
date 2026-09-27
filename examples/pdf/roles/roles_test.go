package roles

import "testing"

func collectRoleWords(source string) map[string]Role {
	words := map[string]Role{}
	for _, span := range TagGoRoles(source) {
		words[source[span.Start:span.End]] = span.Role
	}

	return words
}

func TestTagGoRolesNamesEntitiesAndActions(t *testing.T) {
	source := `type Student struct{ Name FullName }

func RenderStudentView(s Student) wire.StudentView {
	return wire.StudentView{Name: strings.TrimSpace(string(s.Name))}
}`
	words := collectRoleWords(source)

	cases := map[string]Role{
		"Student":           Entity,
		"FullName":          Entity,
		"StudentView":       Entity,
		"RenderStudentView": Action,
		"TrimSpace":         Action,
	}
	for word, wanted := range cases {
		if words[word] != wanted {
			t.Errorf("%s wanted Role %d, got %d", word, wanted, words[word])
		}
	}
}

// A Statement alone is not a File; it Stays Plain rather than Guessed.
func TestTagGoRolesLeavesAFragmentPlain(t *testing.T) {
	if spans := TagGoRoles(`return nil, err`); spans != nil {
		t.Fatalf("a Fragment Must stay Plain, got %v", spans)
	}
}
