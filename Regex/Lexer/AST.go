package lexer

var _ parser = (*AST)(nil)

type AST struct {
	expr    string
	pos     int
	curChar byte
	Root    *Token
}

func NewAST(expr string) *AST {
	ret := &AST{
		expr:    expr,
		pos:     0,
		curChar: '\\',
		Root:    nil,
	}
	ret.toAST()
	return ret
}

/*
	BNF:
	expr ::= term("|"term)*
	term ::= factor*
	factor ::= (subfactor|subfactor ("*"|"+"|"?"))
	subfactor ::= char | "(" expr ")"
*/

type parser interface {
	nextChar() byte
	parse_expr() *Token
	parse_term() *Token
	parse_factor() *Token
	parse_subfactor() *Token
	toAST()
}

/*implement of parser*/
func (ast *AST) nextChar() (ret byte) {
	if ast.pos == len(ast.expr) {
		return '\\'
	}
	ret = ast.expr[ast.pos]
	ast.pos++
	return
}

func (ast *AST) parse_expr() *Token {
	t := ast.parse_term()
	for ast.curChar == '|' {
		ast.curChar = ast.nextChar()
		p := ast.parse_term()
		alt := NewToken(TOKEN_ALTER, '|', t, p)
		t = alt
	}
	return t
}

func (ast *AST) parse_term() *Token {
	t := ast.parse_factor()
	if ast.curChar == '\\' {
		return t
	}
	for ast.curChar == '(' || isAlphabet(ast.curChar) {
		p := ast.parse_factor()
		concat := NewToken(TOKEN_CONCAT, '@', t, p)
		t = concat
	}
	return t
}

func (ast *AST) parse_factor() *Token {
	t := ast.parse_subfactor()
	if ast.curChar == '\\' {
		return t
	}
	for ast.curChar == '*' || ast.curChar == '+' || ast.curChar == '?' {
		tokenType := -1
		if ast.curChar == '*' {
			tokenType = TOKEN_KLEEN
		} else if ast.curChar == '+' {
			tokenType = TOKEN_P_KLEEN
		} else if ast.curChar == '?' {
			tokenType = TOKEN_OPTION
		}

		ast.curChar = ast.nextChar()

		k := &Token{}
		if tokenType == TOKEN_KLEEN {
			k = NewToken(tokenType, '*', t, nil)
		} else if tokenType == TOKEN_P_KLEEN {
			k = NewToken(tokenType, '+', t, nil)
		} else if tokenType == TOKEN_OPTION {
			k = NewToken(tokenType, '?', t, nil)
		}
		t = k
	}
	return t
}

func (ast *AST) parse_subfactor() *Token {
	t := &Token{}
	if ast.curChar == '(' {
		ast.curChar = ast.nextChar()
		t = ast.parse_expr()
		//skip ')'
		ast.curChar = ast.nextChar()
	} else if ast.curChar == ')' {
		ast.curChar = ast.nextChar()
	} else {
		t = NewToken(TOKEN_CHAR, ast.curChar, nil, nil)
		ast.curChar = ast.nextChar()
	}
	return t
}

/*other help functions*/
func (ast *AST) toAST() {
	if ast.expr == "" {
		return
	}
	ast.curChar = ast.nextChar()
	ast.Root = ast.parse_expr()
}
