# Go语言实现子集构造法

在本系列的第二章我们已经详细介绍过子集构造法了，现在我们用一张图做一个简单的回顾

![子集构造法](%E7%94%B1%E5%AE%9E%E4%BE%8B%E5%BC%80%E5%A7%8B%E6%9E%84%E9%80%A0%E6%AD%A3%E5%88%99%E5%BC%95%E6%93%8E(5)Go%E8%AF%AD%E8%A8%80%E5%AE%9E%E7%8E%B0%E5%AD%90%E9%9B%86%E6%9E%84%E9%80%A0%E6%B3%95/%E5%AD%90%E9%9B%86%E6%9E%84%E9%80%A0%E6%B3%95.PNG)

需要声明一下，我们这里也把`Dtran`函数称作`delta`操作（即我们在定义一个状态机如何转移时用的那个数学符号），同样表示从某个状态（或者状态集合）经由某条边（字符）可以到达的状态。再次回顾一下`epsilon`闭包指的是从某个状态（状态集合）出发一直走空边能够到达的所有状态。

## DFA的相关结构

链式前向星已经在第三章介绍过了，这里不多介绍。

### DFA_Node

```golang
type DFA_Node struct {
	Id       int			//每个节点都有唯一的标识符
	IsAccept bool			//当前节点是否是接收状态，不同于NFA只有一个起始一个接受，DFA可以多个接受
	StatesIn map[int]bool	//当前状态的集合，由多个NFA的节点编号组成
	Edges    *DFA_Edge		//核心结构，每一条从当前节点出发的边都会存在这个链表中，存储链表头节点
	Next     *DFA_Node		//用于遍历，在构建过程中存储新节点链条
}
```

### DFA_Edge

```golang
type DFA_Edge struct {
	Char byte
	From *DFA_Node
	To   *DFA_Node
	Next *DFA_Edge
}
```

完全和`NFA_Node`相同

### DFA_Graph

```golang
type DFA_Graph struct {
	innerNFA    *nfa.NFA_Graph		//内建的NFA_Graph结构
	NodeCount   int           		//所有图节点的个数，随着构造过程逐渐增加
	Alphabeta   map[byte]bool 		//所有可行的字符集合
	Head        *DFA_Node     		//整个DFA图的起始节点，构成一条链表用于最后的存图
	Nodes       []*DFA_Node   		//存图后可以通过该项索引所有的点
	Edges       map[byte][]*DFA_Edge//存图后可以通过相关的字符索引对应的边
	NeedNewNode bool				//在子集构造法中标记对象
}
```

## DFA的核心函数

存图相关的内容此处不再赘述。核心函数主要包含了三个，分别是：

1. `delta`
2. `epsilonClosure`
3. ` subsetContruct`

### delta函数实现

在链式前向星结构存图的加成下，`delta`函数的实现十分简单，只需要通过传入的字符索引一个`NFA_Edge`的数组，然后遍历该数组检查边的起点是否在传入的集合中，然后将边的到达节点加入返回集合即可。

```golang
func (g *DFA_Graph) delta(s map[int]bool, ch byte) (ret map[int]bool) {
	edges := g.innerNFA.Edges[ch]
	for _, edge := range edges {
		if s[edge.From.Id] {
			ret[edge.To.Id] = true
		}
	}
	return
}
```

### epsilonClosure实现

`epsilonClosure`的实现稍微复杂一点。我们在`NFA_Node`的结构中安排了`isVisited`这一属性同时还有一个`ResetNodeVisited`函数，这用于给`DFA`实现广度优先搜索。对于传入的集合经由每一条空边能到达的未访问状态就加入队列，直到队列为空，将队列中被加入过的状态标记为访问过。最后把所有节点访问状态清零。

```golang
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
```

### subsetConstruct实现

实现了上述的两个操作后，只要按部就班的照着算法写就可以了，我们用一个`workQueue`来表示`Dstates`，整个过程基本就是翻译伪代码：

```golang
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
```

## 接口规范

```golang
type DFA interface {
	findOrInsertNode(map[int]bool) *DFA_Node		//以集合作为索引查找节点
	addEdge(from, to map[int]bool, ch byte)			//添加一条新边
	delta(s map[int]bool, ch byte) map[int]bool		
	epsilonClosure(x map[int]bool) map[int]bool
	constructAlphabeta()
	subsetContruct()
	storeGraph()									//存图
	ToDFA()											//将NFA转换为DFA
}
```

