package regex

import mdfa "goexpr/Regex/MDFA"

var _ Regex = (*regex)(nil)

type Regex interface {
	Replace(string)
	Match(string) bool
	Search(string, []string) bool
	nextChar() byte
	roolBack()
}

type regex struct {
	innerMDFA *mdfa.MDFA_Graph
	testStr   string
	pos       int
	curChar   byte
	nodes     []*mdfa.MDFA_Node
	matchStr  string
}

func RE(str string) Regex {
	re := &regex{}
	re.pos = 0
	re.innerMDFA = mdfa.NewMDFA_Graph(str)
	re.nodes = re.innerMDFA.Nodes
	return re
}

//TODO
func (re *regex) Replace(string)

func (re *regex) Match(tstr string) bool {
	re.pos = 0
	re.testStr = tstr
	re.curChar = re.nextChar()
	isMatched := false
	state := 0
S:
	if re.nodes[state].IsAccept && re.curChar == '\r' {
		isMatched = true
		return isMatched
	}
	edge := re.nodes[state].Edges
	for edge != nil {
		if edge.Ch == re.curChar {
			re.curChar = re.nextChar()
			state = edge.To.Id
			goto S
		}
		edge = edge.Next
	}
	if edge == nil {
		isMatched = false
	}
	return isMatched
}

//TODO
func (re *regex) Search(tstr string, result []string) bool {
	// result = make([]string, 0)
	// re.pos = 0
	// re.testStr = tstr
	// state := 0
	// stk := make([]int, 0)
	return false
}

func (re *regex) nextChar() (ret byte) {
	if re.pos == len(re.testStr) {
		return '\r'
	}
	ret = re.testStr[re.pos]
	re.pos++
	return
}

func (re *regex) roolBack() {
	re.pos--
}
