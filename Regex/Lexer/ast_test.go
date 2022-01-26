package lexer

import "testing"

func Test_AST(t *testing.T) {
	re0 := "abc*d|(cfuck)?tt"
	re1 := "a(b|c)*(cdef)"

	tree0 := NewAST(re0)
	tree1 := NewAST(re1)

	PrintTree(tree0.Root)
	PrintTree(tree1.Root)
}
