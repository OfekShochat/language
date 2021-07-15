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

type FunctionDecleration struct {
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

func ParseFunction(tokens []Token) (FunctionDecleration, int) {
	f := FunctionDecleration{}
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
			a, _ := Parse(tokens[i : i+length]) // TODO(ghostway): make Parse() also parse Variable Names.
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

func ParseNumberExpressions(tokens []Token) (Node, error, int) {
	if len(tokens) == 1 && tokens[0].Type == 'n' {
		return Node{Value: tokens[0].Value, Type: "Number"}, nil, 1
	}
	f := Node{}
	switch tokens[1].Value {
	case "+":
		f.Value = "plus"
	case "-":
		f.Value = "sub"
	case "*":
		f.Value = "mul"
	default:
		return Node{}, fmt.Errorf("Syntax Error '%s' '%s %s'", GetTokenString(tokens[0:2]), tokens[0].Value, tokens[1].Value), -1
	}
	f.Type = "ConstantMathExpression"

	f.Params = append(f.Params, Node{Value: tokens[0].Value, Type: "Number"})
	i := 2
	for ; i < len(tokens); i++ {
		node, _, n := ParseNumberExpressions(tokens[i:])

		i += n - 1
		f.Params = append(f.Params, node)
		f.Params = append(f.Params)
	}
	return f, nil, i
}

func ParseKeywords(tokens []Token) (Node, FunctionDecleration, error, int) {
	if tokens[0].Type == 'F' {
		f, n := ParseFunction(tokens[1:])
		return Node{}, f, nil, n
	} else if tokens[0].Type == 'R' {
		node, err, n := ParseExpression(tokens[1:])
		returns_node := Node{Type: "Keyword", Value: "Returns"}
		returns_node.Params = append(returns_node.Params, node)
		return returns_node, FunctionDecleration{}, err, n
	}
	return Node{}, FunctionDecleration{}, fmt.Errorf("Didn't find any keywords in %s", GetTokenString(tokens[0:1])), -1
}

func ParseExpression(tokens []Token) (Node, error, int) {
	if tokens[0].Type == 'n' {
		node, err, n := ParseNumberExpressions(tokens)
		if err != nil {
			panic(err)
		}
		return node, err, n
	}
	return Node{}, fmt.Errorf("Didn't find expression in %s", GetTokenString(tokens[0:1])), -1
}

func Parse(tokens []Token) ([]Node, []FunctionDecleration) {
	var result []Node
	var functions []FunctionDecleration
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
	}
	return result, functions
}
