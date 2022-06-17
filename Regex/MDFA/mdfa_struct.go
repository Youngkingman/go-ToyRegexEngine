package mdfa

import dfa "goexpr/Regex/DFA"

type MDFA interface {
	split(map[int]bool) SetPair
	hopcroft()
	addEdge(from, to int, c byte)
	findnode(int) int
	MinimizeDFA()
}

type SetPair struct {
	Set1 map[int]bool
	Set2 map[int]bool
}

type MDFA_Node struct {
	Id       int
	Set      map[int]bool // set points to all vertex in DFA
	IsAccept bool
	Edges    *MDFA_Edge
	Next     *MDFA_Node
}

func NewMDFA_Node(id int, s map[int]bool, isacpt bool, edge *MDFA_Edge, next *MDFA_Node) (ret *MDFA_Node) {
	ret.Id = id
	ret.Set = s
	ret.Edges = edge
	ret.Next = next
	ret.IsAccept = isacpt
	return
}

type MDFA_Edge struct {
	Ch   byte
	From *MDFA_Node
	To   *MDFA_Node
	Next *MDFA_Edge
}

func NewMDFA_Edge(ch byte, from, to *MDFA_Node, next *MDFA_Edge) (ret *MDFA_Edge) {
	ret.Ch = ch
	ret.From, ret.To = from, to
	ret.Next = next
	return
}

type MDFA_Graph struct {
	innerDFA  dfa.DFA_Graph
	NodeCount int
	Head      *MDFA_Node
	Nodes     []*MDFA_Node
	AlphaBeta map[byte]bool
	/*Used to compare whether two divisions are different*/
	GammaOld []map[int]bool
	GammaNew []map[int]bool
}
