package lexer

import "testing"

func Test_AST(t *testing.T) {
	//re1 := "abc*d|(cfuck)?tt"
	re1 := "a(b|c)*(cdef)"
	tree1 := NewAST(re1)
	PrintTree(tree1.Root)
}
