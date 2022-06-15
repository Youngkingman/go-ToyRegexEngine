package dfa

import nfa "goexpr/Regex/NFA"

type DFA interface {
	FindOrInsertNode(map[int]bool)
	AddEdge(from, to map[int]bool, ch byte)
	Delta(s map[int]bool, ch byte)
	Epsilon_Closure(x map[int]bool)
	SubsetContruct()
	StoreGraph()
}

type DFA_Node struct {
	Num      int
	IsAccept bool
	Susbsets map[int]bool
	Edges    *DFA_Edge
	Next     *DFA_Node
}

type DFA_Edge struct {
	Char byte
	From *DFA_Node
	To   *DFA_Node
	Next *DFA_Edge
}

type DFA_Graph struct {
	innerNFA    nfa.NFA
	NodeCount   int
	Alphabeta   map[byte]bool
	Head        *DFA_Node
	Nodes       []*DFA_Node
	Edges       map[byte][]*DFA_Edge
	NeedNewNode bool
}

func NewDFA_Node() *DFA_Node

func NewDFA_Edge() *DFA_Edge
