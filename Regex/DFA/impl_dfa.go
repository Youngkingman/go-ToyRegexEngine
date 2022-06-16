package dfa

import nfa "goexpr/Regex/NFA"

var _ DFA = (*DFA_Graph)(nil)

func (g *DFA_Graph) toDFA() {
	g.ConstructAlphabeta()
	g.SubsetContruct()
	g.StoreGraph()
}

func (g *DFA_Graph) FindOrInsertNode(set map[int]bool) *DFA_Node {
	node := g.Head
	for node != nil {
		if compareSet(set, node.StatesIn) {
			return node
		}
		node = node.Next
	}
	p := NewDFA_Node(g.NodeCount, false, set, nil, nil)
	g.NodeCount++
	if set[g.innerNFA.AcceptId] {
		p.IsAccept = true
	}
	g.NeedNewNode = true
	p.Next = g.Head
	g.Head = p
	return p
}

func (g *DFA_Graph) AddEdge(from, to map[int]bool, ch byte) {
	f := g.FindOrInsertNode(from)
	t := g.FindOrInsertNode(to)
	edge := NewDFA_Edge(ch, f, t, f.Edges)
	f.Edges = edge
}

// This operation returns the set of next states from set `s` through `ch`, which call delta
func (g *DFA_Graph) Delta(s map[int]bool, ch byte) (ret map[int]bool) {
	edges := g.innerNFA.Edges[ch]
	for _, edge := range edges {
		if s[edge.From.Id] {
			ret[edge.To.Id] = true
		}
	}
	return
}

func (g *DFA_Graph) Epsilon_Closure(x map[int]bool) (epsilonClosure map[int]bool) {
	Q := make([]int, 0)
	for k, _ := range x {
		Q = append(Q, k)
		for len(Q) > 0 {
			q := Q[0]
			Q = Q[1:]
			epsilonClosure[q] = true
			p := g.innerNFA.Nodes[q].Edges
			for p != nil && p.Char == nfa.EPSILON_EDGE {
				if !p.To.IsVisited {
					Q = append(Q, p.To.Id)
					p.To.IsVisited = true
				}
				p = p.Next
			}
		}
	}
	g.innerNFA.ResetNodesVisited()
	return
}

func (g *DFA_Graph) ConstructAlphabeta() {
	edges := g.innerNFA.Edges
	for ch, v := range edges {
		if v != nil {
			g.Alphabeta[ch] = true
		}
	}
}

func (g *DFA_Graph) SubsetContruct() {
	initSet := map[int]bool{
		g.innerNFA.StartId: true,
	}
	q0 := g.Epsilon_Closure(initSet)
	workQueue := make([]map[int]bool, 0)
	workQueue = append(workQueue, q0)
	for len(workQueue) > 0 {
		q := workQueue[0]
		workQueue = workQueue[1:]
		for c, _ := range g.Alphabeta {
			v := g.Delta(q, c)
			if v == nil {
				continue
			}
			t := g.Epsilon_Closure(v)
			g.AddEdge(q, t, c)
			if g.NeedNewNode == true {
				workQueue = append(workQueue, t)
				g.NeedNewNode = false
			}
		}
	}
}

func (g *DFA_Graph) StoreGraph() {
	g.Nodes = make([]*DFA_Node, g.NodeCount)
	for p := g.Head; p != nil; p = p.Next {
		g.Nodes[p.Id] = p
		for q := p.Edges; q != nil; q = q.Next {
			g.Edges[q.Char] = append(g.Edges[q.Char], q)
		}
	}
}
