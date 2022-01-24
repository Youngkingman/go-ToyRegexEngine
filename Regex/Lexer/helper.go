package lexer

//判断是否ASCII码和关键字
func isAlphabet(ch byte) bool {
	return (ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

//打印AST树
func PrintTree(root *Token) {
	q := make([]*Token, 0)
	q = append(q, root)
	for len(q) > 0 {
		l := len(q)
		for i := 0; i < l; i++ {
			print(string(q[0].Char), " ")
			if q[0].LeftChild != nil {
				q = append(q, q[0].LeftChild)
			}
			if q[0].RightChild != nil {
				q = append(q, q[0].RightChild)
			}
			q = q[1:]
		}
		println("")
	}
}
