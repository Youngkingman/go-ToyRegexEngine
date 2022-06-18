# Go语言实现Thompson算法

在本系列的第二章我们已经详细介绍过`Thompson`算法了，现在我们用一张图做一个简单的回顾

![](%E7%94%B1%E5%AE%9E%E4%BE%8B%E5%BC%80%E5%A7%8B%E6%9E%84%E9%80%A0%E6%AD%A3%E5%88%99%E5%BC%95%E6%93%8E(4)Go%E8%AF%AD%E8%A8%80%E5%AE%9E%E7%8E%B0Thompson%E7%AE%97%E6%B3%95/Thompson.png)

我们会根据解析的语法树递归地将最终的`NFA`逐步构建出来。

## NFA相关的数据结构（链式前向星简介）

首先在这里给出构造`NFA`所需的数据结构，同时介绍构造图的链式前向星结构，后续的`DFA`以及`MDFA`都会采用这种方式进行图的存储。

链式前向星首先是将所有的节点放在一个数组中，这样我们就能通过其`Id`进行`O(1)`的复杂度索引，空间复杂度为`O(N)`，之后每个节点同时持有所有从该节点出发的边链表，这样查询下一个节点的复杂度就是`O(m)`。这样做的好处：不仅在点数量多和边数量多的情况之间进行了取舍，对边结构的建模使得我们能够存储更多的信息。像网络流之类的题目通常也会采取链式前向星进行建模。

### NFA_Node

```golang
type NFA_Node struct {
	Id        int		//每个节点都有唯一的标识符
	IsVisited bool		//用于DFA构造的广度优先搜索，此处可先忽略
	Next      *NFA_Node	//用于遍历，在构建过程中存储新节点链条
	Edges     *NFA_Edge	//核心结构，每一条从当前节点出发的边都会存在这个链表中，存储链表头节点
}
```

### NFA_Edge

```golang
type NFA_Edge struct {
	Char byte		//该边的持有字符
	From *NFA_Node	//该边的起始节点
	To   *NFA_Node	//该边到达节点
	Next *NFA_Edge // 单向链表的下一个节点
}
```

### NFA_Graph

```golang
type NFA_Graph struct {
	AstTree   *lexer.AST			//通过调用Lexer获取的解析树
	Head      *NFA_Node				//整个NFA图的起始节点，构成一条链表用于最后的存图
	NodeCount int					//所有图节点的个数，随着构造过程逐渐增加
	StartId   int					//当前图起始节点的标识符
	AcceptId  int					//当前图接受（终止）节点的标识符
	Nodes     []*NFA_Node			//存图后可以通过该项索引所有的点
	Edges     map[byte][]*NFA_Edge	//存图后可以通过相关的字符索引对应的边
}
```

### 与图的构建相关的三个函数

这三个函数在`DFA`和`MDFA`中共同拥有，都是用于链式前向星图的构建，它们分别是`findOrInsertNode(num *int*) *NFA_Node`|`addEdge(from, to *int*, edgeChar *byte*)`|`storeGraph() `。

```golang
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
```

`findOrInsertNode`函数在每次创建图节点的时候先进行查找，如果节点存在则直接返回，否则创建新节点并将其加入链表。

```golang
func (g *NFA_Graph) addEdge(from, to int, edgeChar byte) {
	fNode := g.findOrInsertNode(from)
	tNode := g.findOrInsertNode(to)
	edge := NewNFA_Edge(edgeChar, fNode, tNode, fNode.Edges)
	fNode.Edges = edge
}
```

`addEdge`用于在两个节点之间构建一条新的边。

```golang
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
```

`storeGraph`用于

- 将已经构建的节点链表整合成数组，用于下一步的快速索引
- 将每一个节点所带的边链表加入集合索引，这样我们可以通过字符索引所有相关的边

## 核心Thompson算法的实现

`Thompson`算法从根节点开始，根据每个节点的类型决定接下来的行为：

1. 如果当前节点是字符即`TOKEN_CHAR`，则直接创建两个节点和链接两个节点的边，图的起始和接收节点设为当前两个节点
2. 如果当前节点是`|`即`TOKEN_ALTHER`，则
   1. 递归左子树，存储返回图的起始和接受节点
   2. 再递归右子树，存储返回图的起始和接受节点
   3. 建立两个节点作为新的起始节点和接受节点，通过并联方式建立新的图（添加空边）
3. 如果当前节点是连接即`TOKEN_CONCAT`，则
   1. 递归左子树，存储返回图的起始和接受节点
   2. 再递归右子树，存储返回图的起始和接受节点
   3. 在左子图的接受节点和右子图的起始节点间建立空边
   4. 图的起始节点设置为左子图的起点，接受节点设置为右子图的接受节点
4. 如果当前节点是`*`即`TOKEN_KLEEN`，则
   1. 递归左子树，存储返回图的起始和接受节点
   2. 在左子图的**接受点**和**起始点**间添加一条空边（注意反向）
   3. 建立两个节点作为新的起始节点和接受节点，之间建立一条空边，同时通过空边和之前的子图串联
5. 如果当前节点是`?`即`TOKEN_OPTION`，则类似于克林闭包
   1. 递归左子树，存储返回图的起始和接受节点
   2. 建立两个节点作为新的起始节点和接受节点，在起始和接受节点添加一条空边，同时通过空边和之前的子图串联
6. 如果当前节点是`+`即至少重复一次，则
   1. 递归左子树
   2. 在图的**接收点**和**起始点**间添加一条空边（注意反向）

上述的说明可以配合代码文件中的`Regex/NFA/impl_nfa.go`中的` recurciveConstruct`函数进行理解。

## 接口规范

```golang
type NFA interface {
	thompsonAlgor(tree *lexer.AST)			//核心算法
	addEdge(from, to int, edgeChar byte)	//添加边
	findOrInsertNode(num int) *NFA_Node		//插入或查询节点
	storeGraph()							//存图
	ToNFA()									//转换为NFA
	ResetNodesVisited()						//用于DFA构建时的广度优先搜索，此处可先忽略
}
```

