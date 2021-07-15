package main

import (
	"fmt"
)

type Node struct {
	Value  string
	Type   string
	Params []Node
}

type Param struct {
	Name  string
	Type  string
	Value string
}

type FunctionDeclaration struct {
	Name    string
	Params  []Param
	Body    []Node
	Returns string
}

func ParseFunctionParams(tokens []Token) ([]Param, int) {
	var params []Param
	p := Param{}
	i := 0
	for ; i < len(tokens); i++ {
		if tokens[i].Type == 'T' {
			if p.Name == "" {
				p.Type = tokens[i].Value
			}
		} else if tokens[i].Type == 'N' {
			if p.Type != "" {
				p.Name = tokens[i].Value
			} else {
				panic(fmt.Errorf("Name declared before type"))
			}
		} else if tokens[i].Value == "," {
			params = append(params, p)
			p = Param{}
		} else if tokens[i].Value == ")" {
			if p.Type != "" || p.Name != "" {
				params = append(params, p)
			}
			break
		}
	}

	return params, i
}

func GetFunctionBodyLength(tokens []Token) (int, error) {
	for i := 0; i < len(tokens); i++ {
		if tokens[i].Value == "}" {
			return i, nil
		}
	}
	return -1, fmt.Errorf("Didn't find end of function")
}

func ParseFunction(tokens []Token) (FunctionDeclaration, int) {
	f := FunctionDeclaration{}
	f.Name = tokens[0].Value
	i := 1
	for ; i < len(tokens); i++ {
		if tokens[i].Value == "(" {
			ps, l := ParseFunctionParams(tokens[i:])
			f.Params = ps
			i += l
		} else if tokens[i].Value == "{" {
			length, err := GetFunctionBodyLength(tokens[i:])
			if err != nil {
				fmt.Println(err)
			}
			offset := 1
			if tokens[i+1].Type == 'l' {
				offset = 2
			}
			a, _ := Parse(tokens[i+offset : i+length])
			f.Body = a
		} else if tokens[i].Value == "->" {
			f.Returns = tokens[i+1].Value
		}
	}
	//p := Param{Name: tokens[3].Value, Type: tokens[2].Value}
	//f.Params = append(f.Params, p)

	return f, i
}

func GetTokenString(tokens []Token) string {
	result_string := ""
	for i := 0; i < len(tokens); i++ {
		result_string += string(tokens[i].Type)
		if tokens[i].Type == 'c' {
			result_string += tokens[i].Value
		}
		result_string += " "
	}
	return result_string[:len(result_string)-1]
}

func GetExpressionLength(tokens []Token) int {
	i := 0
	used_paren := false
	for ; i < len(tokens); i++ {
		if tokens[i].Value == "(" {
			used_paren = true
		} else if tokens[i].Type == 'l' {
			if used_paren {
				continue
			} else {
				break
			}
		} else if tokens[i].Value == ")" {
			break
		}
	}
	if used_paren {
		return i - 1
	}
	return i
}

func GetExpressionType(node Node) string {
	for i := 0; i < len(node.Params); i++ {

	}
	return "ConstantMathExpression"
}

func ParseNumberExpressions(tokens []Token) (Node, error, int) {
	if len(tokens) == 1 {
		if tokens[0].Type == 'n' {
			return Node{Value: tokens[0].Value, Type: "Number"}, nil, 1
		} else if tokens[0].Type == 'N' {
			return Node{Value: tokens[0].Value, Type: "GetVariable"}, nil, 1
		}
	}
	f := Node{}
	switch tokens[1].Value {
	case "+":
		f.Value = "plus"
	case "-":
		f.Value = "sub"
	case "*":
		f.Value = "mul"
	case "/":
		f.Value = "div"
	default:
		return Node{}, fmt.Errorf("Syntax Error Token String '%s' String '%s %s'", GetTokenString(tokens[0:2]), tokens[0].Value, tokens[1].Value), -1
	}
	//f.Type = "ConstantMathExpression" // TODO(ghostway): check if its really a compile time expression

	if tokens[0].Type == 'n' {
		f.Params = append(f.Params, Node{Value: tokens[0].Value, Type: "Number"})
	} else if tokens[0].Type == 'N' {
		f.Params = append(f.Params, Node{Value: tokens[0].Value, Type: "GetVariable"})
	}
	i := 2
	for ; i < len(tokens); i++ {
		node, _, n := ParseNumberExpressions(tokens[i:])

		i += n - 1
		f.Params = append(f.Params, node)
		f.Params = append(f.Params)
	}
	f.Type = GetExpressionType(f)
	return f, nil, i
}

func ParseVariable(tokens []Token, t byte) (Node, error, int) {
	v := Node{}
	if t == 't' {
		v.Type = "CompileTimeVariableDeclaration"
	} else {
		v.Type = "VariableDeclaration"
	}
	a := 3
	if t == 't' {
		v.Value = tokens[1].Value
		a = 4
	} else {
		v.Value = tokens[0].Value

	}
	params, err, n := ParseExpression(tokens[a : a+GetExpressionLength(tokens[a:])])
	v.Params = append(v.Params, params)
	return v, err, n + a
}

func ParseKeywords(tokens []Token) (Node, FunctionDeclaration, error, int) {
	if tokens[0].Type == 'F' {
		f, n := ParseFunction(tokens[1:])
		return Node{}, f, nil, n
	} else if tokens[0].Type == 'R' {
		node, err, n := ParseExpression(tokens[1:])
		returns_node := Node{Type: "Keyword", Value: "Returns"}
		returns_node.Params = append(returns_node.Params, node)
		return returns_node, FunctionDeclaration{}, err, n
	} else if tokens[0].Type == 'T' || tokens[0].Type == 't' {
		node, err, n := ParseVariable(tokens, tokens[0].Type)
		if err != nil {
			panic(err)
		}
		return node, FunctionDeclaration{}, nil, n
	}
	return Node{}, FunctionDeclaration{}, fmt.Errorf("Didn't find any keywords in %s", GetTokenString(tokens[0:1])), -1
}

func ParseExpression(tokens []Token) (Node, error, int) {
	if tokens[0].Type == 'n' || tokens[0].Type == 'N' {
		node, err, n := ParseNumberExpressions(tokens[:GetExpressionLength(tokens)])
		if err != nil {
			panic(err)
		}
		return node, err, n
	}
	return Node{}, fmt.Errorf("Didn't find expression in %s", GetTokenString(tokens[0:1])), -1
}

func Parse(tokens []Token) ([]Node, []FunctionDeclaration) {
	var result []Node
	var functions []FunctionDeclaration
	for i := 0; i < len(tokens); i++ {
		node, function, err, n := ParseKeywords(tokens[i:])
		if err == nil {
			i += n
			if node.Type == "" {
				functions = append(functions, function)
			} else {
				result = append(result, node)
			}
			continue
		}
		node, err, n = ParseExpression(tokens[i:])
		if err == nil {
			result = append(result, node)
			i += n
			continue
		}
		if tokens[i].Type == 'l' {
			i++
			continue
		}
		panic(fmt.Sprintf("unexpected token type %s (%s)", string(tokens[i].Type), tokens[i].Value))
	}
	return result, functions
}
