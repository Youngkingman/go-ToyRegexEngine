<center><h1>Go语言实现递归下降分析
## 项目基本介绍


现在我们正式开始基于Go语言构建一个正则表达式引擎。项目地址位于[这里](https://github.com/Youngkingman/go-ToyRegexEngine)。

在这之前我们先整体介绍一下项目的基本架构。正则引擎位于`Regex`文件夹中，下面包含了四个文件夹`Lexer`|`NFA`|`DFA`|`MDFA`，分别是**语法解析器**/**非确定性有限状态自动机**/**确定性有限状态自动机**和**最小化确定性有限状态自动机**。

我们向正则引擎传递一个正则表达式`regexpr`字符串，接下来：

1. 这个字符串会首先传送给`MDFA`，`MDFA`调用`DFA`的构造方法将字符串传递正则字符串
2. `DFA`接受`regexpr`后调用`NFA`的构造方法传递正则字符串
3. `NFA`接受`regexpr`后调用`Lexer`的构造方法传递正则字符串
4. `Lexer`接受`regexpr`生成语法解析树`AST`，并回传给`NFA`
5. `NFA`根据回传的内建`AST`经由`Thompson`算法构建`NFA_Graph`并回传给`DFA`
6. `DFA`根据回传的内建`NFA_Graph`经由子集构造法构建`DFA_Graph`并回传给`MDFA`
7. `MDFA`根据回传的`DFA_Graph`经由`Hopcroft`算法构建`MDFA_Graph`并回传给`Regex`，正则引擎通过该最小化的确定性有限状态自动机进行字符串匹配

在这个过程中向下传递的过程通过每个结构各自的构造函数实现，而回传过程通过`Golang`的`interface`规范了相互之间调用的接口，同时每一步都采用了**空接口断言的方式**保证了所有接口都有被实现。第5~7步中的`*_Graph`采用链式前向星的数据结构进行图的存储，在邻接表和邻接矩阵两种极端的结构之间取了一个中庸性能的结构，因为在算法实现过程中我们需要大量访问边和点而并非只是其中之一，所以这可能是一个更好的选择。

## 语法解析树的构建

这部分的工作属于编译器前端的一部分，主要是将文本根据给定的巴克斯范式解析为一系列`token`，生成一颗语法解析树，之前我们已经给出过BNF范式的表达如下：

```
BNF:
	expr ::= term("|"term)*
	term ::= factor*
	factor ::= (subfactor|subfactor ("*"|"+"|"?"))
	subfactor ::= char | "(" expr ")"
```

我们可选的工具很多比如`yacc`/`gyacc`等等，但是由于正则表达式属于最简单的语言了，在这个项目里就直接自己手写递归下降解析为`token`了。

### Token结构

```golang
type Token struct {
	Token_Type int
	Char       byte
	LeftChild  *Token
	RightChild *Token
}
```

`Token`是一个典型的二叉树结构，语法解析树是由一个个`Token`所构成的，`Char`字段表示其在正则表达式中的字符（或单纯的连接符号）,`Token_Type`字段由枚举结构给出，总共有六种：

```golang
const (
	TOKEN_CHAR = iota	//字符类型
	TOKEN_ALTER			//或类型`|`
	TOKEN_CONCAT		//链接符
	TOKEN_KLEEN			//克林闭包`*`
	TOKEN_OPTION		//可选类型`?`
	TOKEN_P_KLEEN		//重复一次或以上`+`
)
```

### AST结构以及其对parser接口的实现

```golang
type AST struct {
	expr    string
	pos     int
	curChar byte
	Root    *Token
}
```

`AST`的结构包括了内建的正则表达式`expr`，当前解析的位置`pos`，当前解析的字符`curChar`以及当前解析树的根节点`Root`。

`parser`接口的定义如下：

```golang
type parser interface {
	nextChar() byte			//解析下一个字符
	parse_expr() *Token		//对应BNF中的第一行
	parse_term() *Token		//对应BNF中的第二行
	parse_factor() *Token	//对应BNF中的第三行
	parse_subfactor() *Token//对应BNF中的第四行
	toAST()					//进行字符串到AST的转换
}
```

该接口将由`AST`结构实现

具体代码直接参考项目源文件[这里](https://github.com/Youngkingman/go-ToyRegexEngine)。
