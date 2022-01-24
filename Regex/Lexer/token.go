package lexer

const (
	TOKEN_CHAR = iota
	TOKEN_ALTER
	TOKEN_CONCAT
	TOKEN_KLEEN
	TOKEN_OPTION
	TOKEN_P_KLEEN
)

type Token struct {
	Token_Type int
	Char       byte
	LeftChild  *Token
	RightChild *Token
}

func NewToken(Token_Type int, Char byte, LeftChild, RightChild *Token) *Token {
	ret := Token{
		Token_Type: Token_Type,
		LeftChild:  LeftChild,
		RightChild: RightChild,
		Char:       Char,
	}
	return &ret
}
