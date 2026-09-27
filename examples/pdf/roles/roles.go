// Package roles Colours Go Code by what each Name Does in the Story,
// not by what the Syntax Calls it. A Type is an Entity, a Call is an
// Action. Go's own AST Answers which is which — no Judgment, a Parse.
package roles

import (
	"go/ast"
	"go/parser"
	"go/token"
	"sort"
)

// Role Names the Part a Word Plays.
type Role int

const (
	Plain Role = iota
	Entity
	Action
)

// Span Marks one Word by its Byte Offsets in the Source.
type Span struct {
	Start, End int
	Role       Role
}

// fragmentHeader Lets a Snippet without a Package Clause still Parse.
const fragmentHeader = "package fragment\n"

// TagGoRoles Finds every Entity and Action in a Go Snippet, in Order.
// A Snippet that is not a whole Declaration Returns nothing:
// an honest Plain Block Beats a guessed Colour.
func TagGoRoles(source string) []Span {
	file, err := parser.ParseFile(token.NewFileSet(), "", fragmentHeader+source, parser.SkipObjectResolution)
	if err != nil {
		return nil
	}

	found := map[int]Span{}
	mark := func(name *ast.Ident, role Role) {
		if name == nil || name.Name == "_" {
			return
		}
		start := int(name.Pos()) - 1 - len(fragmentHeader)
		found[start] = Span{Start: start, End: start + len(name.Name), Role: role}
	}

	ast.Inspect(file, func(node ast.Node) bool {
		switch v := node.(type) {
		case *ast.TypeSpec:
			mark(v.Name, Entity)
		case *ast.Field:
			markTypeNames(v.Type, mark)
		case *ast.CompositeLit:
			markTypeNames(v.Type, mark)
		case *ast.ValueSpec:
			markTypeNames(v.Type, mark)
		case *ast.FuncDecl:
			mark(v.Name, Action)
		case *ast.CallExpr:
			mark(calledName(v.Fun), Action)
		}

		return true
	})

	return sortSpans(found)
}

// markTypeNames Walks a Type Expression down to the Names inside it.
func markTypeNames(expr ast.Expr, mark func(*ast.Ident, Role)) {
	switch v := expr.(type) {
	case *ast.Ident:
		mark(v, Entity)
	case *ast.SelectorExpr:
		mark(v.Sel, Entity)
	case *ast.StarExpr:
		markTypeNames(v.X, mark)
	case *ast.ArrayType:
		markTypeNames(v.Elt, mark)
	case *ast.MapType:
		markTypeNames(v.Key, mark)
		markTypeNames(v.Value, mark)
	case *ast.IndexExpr:
		markTypeNames(v.X, mark)
		markTypeNames(v.Index, mark)
	}
}

// calledName Returns the Word that Names the Action in a Call.
func calledName(fun ast.Expr) *ast.Ident {
	switch v := fun.(type) {
	case *ast.Ident:
		return v
	case *ast.SelectorExpr:
		return v.Sel
	case *ast.IndexExpr:
		return calledName(v.X)
	}

	return nil
}

func sortSpans(found map[int]Span) []Span {
	spans := make([]Span, 0, len(found))
	for _, span := range found {
		spans = append(spans, span)
	}
	sort.Slice(spans, func(i, j int) bool { return spans[i].Start < spans[j].Start })

	return spans
}
