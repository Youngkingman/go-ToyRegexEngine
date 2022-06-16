package dfa

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

func (g *DFA_Graph) Delta(s map[int]bool, ch byte) map[int]bool

func (g *DFA_Graph) Epsilon_Closure(x map[int]bool) map[int]bool

func (g *DFA_Graph) ConstructAlphabeta() {}

func (g *DFA_Graph) SubsetContruct()

func (g *DFA_Graph) StoreGraph()
