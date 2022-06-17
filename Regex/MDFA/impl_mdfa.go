package mdfa

func (g *MDFA_Graph) split(map[int]bool) SetPair
func (g *MDFA_Graph) hopcroft()
func (g *MDFA_Graph) addEdge(from, to int, c byte)
func (g *MDFA_Graph) findnode(int) int
func (g *MDFA_Graph) MinimizeDFA()
