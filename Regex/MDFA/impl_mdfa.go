package mdfa

// 核心实现，如何进行两个集合的划分
func (g *MDFA_Graph) split(Set map[int]bool) SetPair {
	tmp := copySet(Set)
	s1 := make(map[int]bool)
	s2 := make(map[int]bool)
	for ch, _ := range g.AlphaBeta {
		edges := g.innerDFA.Edges[ch]
		for _, edge := range edges {
			b1 := false
			firstSet := g.GammaOld[0]
			for _, v := range g.GammaOld {
				from := edge.From.Id
				to := edge.To.Id
				if v[to] && !compareSet(tmp, v) {
					if b1 == false {
						firstSet = v
						s1[from] = true
						delete(Set, from)
						b1 = true
					} else {
						if compareSet(firstSet, v) {
							s1[from] = true
							delete(Set, from)
						} else {
							s2[from] = true
							delete(Set, from)
						}
					}
				} else if v[to] && compareSet(tmp, v) {
					s2[from] = true
					delete(Set, from)
				}
			}
		}
		if compareSet(tmp, Set) {
			if len(s1) > 0 {
				s2 = copySet(Set)
				return SetPair{s1, s2}
			} else {
				return SetPair{s2, Set}
			}
		}
	}
	return SetPair{Set1: Set, Set2: Set}
}

func (g *MDFA_Graph) hopcroft() {
	acceptSet := make(map[int]bool)
	notAcceptSet := make(map[int]bool)
	for p := g.innerDFA.Head; p != nil; p = p.Next {
		if p.IsAccept {
			acceptSet[p.Id] = true
		} else {
			notAcceptSet[p.Id] = true
		}
	}
	g.GammaNew = append(g.GammaNew, acceptSet)
	g.GammaNew = append(g.GammaNew, notAcceptSet)
	for !compareTwoGamma(g.GammaOld, g.GammaNew) {
		g.GammaOld = g.GammaNew
		g.GammaNew = nil
	}
	for _, v := range g.GammaOld {
		tmp := v
		s1 := g.split(tmp).Set1
		if len(s1) != 0 {
			g.GammaNew = append(g.GammaNew, s1)
		}
		s2 := g.split(tmp).Set2
		if len(s2) != 0 {
			g.GammaNew = append(g.GammaNew, s2)
		}
	}
}

func (g *MDFA_Graph) addEdge(from, to int, c byte) {
	f := g.findnode(from)
	t := g.findnode(to)
	it := g.Nodes[f].Edges
	for it != nil {
		if it.Ch == c {
			return
		}
		it = it.Next
	}
	p := NewMDFA_Edge(c, g.Nodes[f], g.Nodes[t], g.Nodes[f].Edges)
	g.Nodes[f].Edges = p
}

func (g *MDFA_Graph) findnode(n int) int {
	for i := 0; i < g.NodeCount; i++ {
		if g.Nodes[i].Set[n] {
			return i
		}
	}
	return -1
}

func (g *MDFA_Graph) MinimizeDFA() {
	g.AlphaBeta = g.innerDFA.Alphabeta
	g.hopcroft()

	// store the graph
	for _, v := range g.GammaNew {
		node := make(map[int]bool)
		for k, v1 := range v {
			node[k] = v1
		}
		p := NewMDFA_Node(g.NodeCount, node, false, nil, g.Head)
		g.NodeCount++
		g.Head = p
		for k := range v {
			if g.innerDFA.Nodes[k].IsAccept {
				p.IsAccept = true
			}
		}
	}

	g.Nodes = make([]*MDFA_Node, g.NodeCount)
	for p := g.Head; p != nil; p = p.Next {
		g.Nodes[p.Id] = p
	}

	for p := g.innerDFA.Head; p != nil; p = p.Next {
		for q := p.Edges; q != nil; q = q.Next {
			g.addEdge(q.From.Id, q.To.Id, q.Char)
		}
	}
}
