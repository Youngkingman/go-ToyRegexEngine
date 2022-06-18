package nfa

import lexer "goexpr/Regex/Lexer"

var _ NFA = (*NFA_Graph)(nil)

func (g *NFA_Graph) ToNFA() {
	g.thompsonAlgor(g.AstTree)
	g.storeGraph()
}

/*
	ThompsonAlgor(tree *lexer.AST)
	addEdge(from, to int, edgeChar byte)
	FindOrInsertNode(num int) *NFA_Node
	StoreGraph()
*/

func (g *NFA_Graph) thompsonAlgor(tree *lexer.AST) {
	g.recurciveConstruct(tree.Root)
}

func (g *NFA_Graph) recurciveConstruct(t *lexer.Token) {
	if t.Token_Type == lexer.TOKEN_CHAR {
		/*
			from --ch --> to
		*/
		from := g.NodeCount
		to := g.NodeCount + 1
		g.NodeCount += 2
		g.addEdge(from, to, t.Char)
		g.StartId = from
		g.AcceptId = to
	} else if t.Token_Type == lexer.TOKEN_ALTER {
		/*
					start0 -- (left_sub_graph) --> accept0
				 /											\
			from 											 to
				 \											/
					start1 -- (right_sub_graph) --> accept1
		*/
		g.recurciveConstruct(t.LeftChild)
		start0, acecept0 := g.StartId, g.AcceptId

		g.recurciveConstruct(t.RightChild)
		start1, acecept1 := g.StartId, g.AcceptId

		from, to := g.NodeCount, g.NodeCount+1
		g.NodeCount += 2

		g.addEdge(from, start0, EPSILON_EDGE)
		g.addEdge(from, start1, EPSILON_EDGE)

		g.addEdge(acecept0, to, EPSILON_EDGE)
		g.addEdge(acecept1, to, EPSILON_EDGE)
		g.StartId, g.AcceptId = from, to

	} else if t.Token_Type == lexer.TOKEN_CONCAT {
		/*
			from --ε-->preStart-->(left_sub_graph)-->preAccept--ε-->(right_sub_graph)--ε-->to
		*/
		g.recurciveConstruct(t.LeftChild)
		preStart, preAccept := g.StartId, g.AcceptId
		g.recurciveConstruct(t.RightChild)
		g.addEdge(preAccept, preStart, EPSILON_EDGE)
		g.StartId = preStart
	} else if t.Token_Type == lexer.TOKEN_KLEEN {
		/*
					 -------->--------------ε-------------------->---
					/												 \
				from --ε-->prtStart--(left_sub_graph)--preAccept--ε-->to
			  					\--<-------ε-----------/
		*/
		g.recurciveConstruct(t.LeftChild)
		g.addEdge(g.AcceptId, g.StartId, EPSILON_EDGE)

		from, to := g.NodeCount, g.NodeCount+1
		g.NodeCount += 2
		g.addEdge(from, g.StartId, EPSILON_EDGE)
		g.addEdge(from, to, EPSILON_EDGE)
		g.addEdge(g.AcceptId, to, EPSILON_EDGE)

		g.StartId, g.AcceptId = from, to
	} else if t.Token_Type == lexer.TOKEN_OPTION {
		/*
			     -------->--------------ε-------------------->---
				/												 \
			from --ε-->prtStart--(left_sub_graph)--preAccept--ε-->to
		*/
		g.recurciveConstruct(t.LeftChild)
		from, to := g.NodeCount, g.NodeCount+1
		g.NodeCount += 2

		g.addEdge(from, g.StartId, EPSILON_EDGE)
		g.addEdge(from, to, EPSILON_EDGE)
		g.addEdge(g.AcceptId, to, EPSILON_EDGE)

		g.StartId, g.AcceptId = from, to
	} else if t.Token_Type == lexer.TOKEN_P_KLEEN {
		/*
			from --ε-->prtStart--(left_sub_graph)--preAccept--ε-->to
			  			  \--<-------ε----------------/
		*/
		g.recurciveConstruct(t.LeftChild)
		g.addEdge(g.AcceptId, g.StartId, EPSILON_EDGE)
	}
}

func (g *NFA_Graph) addEdge(from, to int, edgeChar byte) {
	fNode := g.findOrInsertNode(from)
	tNode := g.findOrInsertNode(to)
	edge := NewNFA_Edge(edgeChar, fNode, tNode, fNode.Edges)
	fNode.Edges = edge
}

func (g *NFA_Graph) findOrInsertNode(num int) *NFA_Node {
	node := g.Head
	for node != nil {
		if node.Id == num {
			return node
		}
		node = node.Next
	}
	pNode := NewNFA_Node(num)
	pNode.Next = g.Head
	g.Head = pNode
	return pNode
}

func (g *NFA_Graph) storeGraph() {
	g.Nodes = make([]*NFA_Node, g.NodeCount)
	for p := g.Head; p != nil; p = p.Next {
		g.Nodes[p.Id] = p
		q := p.Edges
		for q != nil && q.Char != EPSILON_EDGE {
			if g.Edges[q.Char] == nil {
				g.Edges[q.Char] = make([]*NFA_Edge, 0)
			}
			g.Edges[q.Char] = append(g.Edges[q.Char], q)
			q = q.Next
		}
	}
}

func (g *NFA_Graph) ResetNodesVisited() {
	p := g.Head
	for p != nil {
		p.IsVisited = false
		p = p.Next
	}
}
