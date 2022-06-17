package dfa

import nfa "goexpr/Regex/NFA"

type DFA interface {
	findOrInsertNode(map[int]bool) *DFA_Node
	addEdge(from, to map[int]bool, ch byte)
	delta(s map[int]bool, ch byte) map[int]bool
	epsilonClosure(x map[int]bool) map[int]bool
	constructAlphabeta()
	subsetContruct()
	storeGraph()
	ToDFA()
}

type DFA_Node struct {
	Id       int
	IsAccept bool
	StatesIn map[int]bool
	Edges    *DFA_Edge
	Next     *DFA_Node
}

type DFA_Edge struct {
	Char byte
	From *DFA_Node
	To   *DFA_Node
	Next *DFA_Edge
}

// Still our graph store in the struct of `链式前向星`
type DFA_Graph struct {
	innerNFA    *nfa.NFA_Graph
	NodeCount   int           //number of vertex
	Alphabeta   map[byte]bool //store all charcters
	Head        *DFA_Node     //current start head
	Nodes       []*DFA_Node   //the same as NFA
	Edges       map[byte][]*DFA_Edge
	NeedNewNode bool
}

func NewDFA_Graph(str string) (ret *DFA_Graph) {
	ret.innerNFA = nfa.NewNFA_Graph(str)
	ret.Head = nil
	ret.NodeCount = 0
	ret.ToDFA()
	return
}

func NewDFA_Node(id int, isAcpt bool, sets map[int]bool, edges *DFA_Edge, next *DFA_Node) *DFA_Node {
	return &DFA_Node{
		Id:       id,
		IsAccept: isAcpt,
		StatesIn: sets,
		Edges:    edges,
		Next:     next,
	}
}

func NewDFA_Edge(ch byte, from *DFA_Node, to *DFA_Node, next *DFA_Edge) *DFA_Edge {
	return &DFA_Edge{
		Char: ch,
		From: from,
		To:   to,
		Next: next,
	}
}
