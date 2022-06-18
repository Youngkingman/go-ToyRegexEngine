package dfa

import nfa "goexpr/Regex/NFA"

var _ DFA = (*DFA_Graph)(nil)

func (g *DFA_Graph) ToDFA() {
	g.constructAlphabeta()
	g.subsetContruct()
	g.storeGraph()
}

func (g *DFA_Graph) findOrInsertNode(set map[int]bool) *DFA_Node {
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

func (g *DFA_Graph) addEdge(from, to map[int]bool, ch byte) {
	f := g.findOrInsertNode(from)
	t := g.findOrInsertNode(to)
	edge := NewDFA_Edge(ch, f, t, f.Edges)
	f.Edges = edge
}

// This operation returns the set of next states from set `s` through `ch`, which call delta
func (g *DFA_Graph) delta(s map[int]bool, ch byte) (ret map[int]bool) {
	edges := g.innerNFA.Edges[ch]
	for _, edge := range edges {
		if s[edge.From.Id] {
			ret[edge.To.Id] = true
		}
	}
	return
}

func (g *DFA_Graph) epsilonClosure(x map[int]bool) (epsilonClosure map[int]bool) {
	Q := make([]int, 0)
	for k := range x {
		Q = append(Q, k)

	}
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
	g.innerNFA.ResetNodesVisited()
	return
}

func (g *DFA_Graph) constructAlphabeta() {
	edges := g.innerNFA.Edges
	for ch, v := range edges {
		if v != nil {
			g.Alphabeta[ch] = true
		}
	}
}

func (g *DFA_Graph) subsetContruct() {
	initSet := map[int]bool{
		g.innerNFA.StartId: true,
	}
	q0 := g.epsilonClosure(initSet)
	workQueue := make([]map[int]bool, 0)
	workQueue = append(workQueue, q0)
	for len(workQueue) > 0 {
		q := workQueue[0]
		workQueue = workQueue[1:]
		for c := range g.Alphabeta {
			v := g.delta(q, c)
			if v == nil {
				continue
			}
			t := g.epsilonClosure(v)
			g.addEdge(q, t, c)
			if g.NeedNewNode {
				workQueue = append(workQueue, t)
				g.NeedNewNode = false
			}
		}
	}
}

func (g *DFA_Graph) storeGraph() {
	g.Nodes = make([]*DFA_Node, g.NodeCount)
	for p := g.Head; p != nil; p = p.Next {
		g.Nodes[p.Id] = p
		for q := p.Edges; q != nil; q = q.Next {
			g.Edges[q.Char] = append(g.Edges[q.Char], q)
		}
	}
}
