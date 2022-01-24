package main

import lexer "goexpr/Regex/Lexer"

func main() {
	re1 := "abc?d+|cfuck*"
	tree1 := lexer.NewAST(re1)
	lexer.PrintTree(tree1.Root)
}
