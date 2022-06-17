package nfa

import lexer "goexpr/Regex/Lexer"

type NFA interface {
	thompsonAlgor(tree *lexer.AST)
	addEdge(from, to int, edgeChar byte)
	findOrInsertNode(num int) *NFA_Node
	storeGraph()
	ToNFA()
	ResetNodesVisited()
}

const EPSILON_EDGE = byte('-')

type NFA_Node struct {
	Id        int
	IsVisited bool
	Next      *NFA_Node
	Edges     *NFA_Edge
}

type NFA_Edge struct {
	Char byte
	From *NFA_Node
	To   *NFA_Node
	Next *NFA_Edge // the property point to last edge with same `from`
}

// The graph is store in the struct of `链式前向星`
type NFA_Graph struct {
	AstTree   *lexer.AST
	Head      *NFA_Node
	NodeCount int
	StartId   int
	AcceptId  int
	Nodes     []*NFA_Node
	Edges     map[byte][]*NFA_Edge
}

func NewNFA_Graph(str string) *NFA_Graph {
	tree := lexer.NewAST(str)
	ret := &NFA_Graph{
		AstTree:  tree,
		Head:     nil,
		StartId:  0,
		AcceptId: 0,
		Nodes:    nil, //make after the construction of nfa graph
		Edges:    make(map[byte][]*NFA_Edge),
	}
	ret.ToNFA()
	return ret
}

func NewNFA_Node(id int) *NFA_Node {
	ret := &NFA_Node{
		Id:        id,
		IsVisited: false,
		Next:      nil,
		Edges:     nil,
	}
	return ret
}

func NewNFA_Edge(ch byte, from, to *NFA_Node, next *NFA_Edge) *NFA_Edge {
	ret := &NFA_Edge{
		Char: ch,
		From: from,
		To:   to,
		Next: next,
	}
	return ret
}
