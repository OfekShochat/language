package main

import (
	"fmt"
	"regexp"
)

type Token struct {
	Type  byte
	Value string
}

func TokenizePattern(pattern *regexp.Regexp, input string) (string, error) {
	var value = ""
	i := 0
	if pattern.FindString(string(input[0])) != "" {
		for ; i < len(input); i++ {
			if pattern.FindString(string(input[i])) == "" {
				break
			}
			value += string(input[i])
		}
	}
	return value, nil
}

func TokenizeNumber(input string) (Token, error, int) {
	pattern := regexp.MustCompile("[[:digit:]]")
	result, err := TokenizePattern(pattern, input)
	if err != nil {
		panic("Something bad happened when tokenizing a Number")
	}
	if result == "" {
		return Token{Type: 'u', Value: ""}, fmt.Errorf("Couldn't parse Number from: %s", input), -1
	}
	/*if i < len(input) {
		return [Token{Type: 'n', Value: result}, Token{Type: 'u', Value: input[i:]}]
	}*/

	return Token{Type: 'n', Value: result}, nil, len(result)
}

func TokenizeString(input string) (Token, error, int) {
	closed := false
	value := ""
	if input[0] == '"' {
		for i := 1; i < len(input); i++ {
			if input[i] == '"' {
				closed = true
				break
			}
			value += string(input[i])
		}
	}
	if !closed {
		return Token{Type: 'u', Value: ""}, fmt.Errorf("Couldn't parse String from: %s", input), 0
	}
	return Token{Type: 's', Value: value}, nil, len(value)
}

func SkipWhiteSpace(input string) int {
	pattern := regexp.MustCompile("[^\\S\r\n]")
	result, err := TokenizePattern(pattern, input)
	if err != nil {
		return 0
	}
	return len(result)
}

func TokenizeChar(input string) (Token, error) {
	switch input[0] {
	case '(':
		return Token{Type: 'c', Value: "("}, nil
	case ')':
		return Token{Type: 'c', Value: ")"}, nil
	case '+':
		return Token{Type: 'c', Value: "+"}, nil
	case '-':
		return Token{Type: 'c', Value: "-"}, nil
	case '*':
		return Token{Type: 'c', Value: "*"}, nil
	case '/':
		return Token{Type: 'c', Value: "/"}, nil
	case '{':
		return Token{Type: 'c', Value: "{"}, nil
	case '}':
		return Token{Type: 'c', Value: "}"}, nil
	case ',':
		return Token{Type: 'c', Value: ","}, nil
	case '\n':
		return Token{Type: 'l', Value: "\n"}, nil
	case ';':
		return Token{Type: 'l', Value: ";"}, nil
	case '=':
		return Token{Type: 'e', Value: "="}, nil
	default:
		return Token{Type: 'u', Value: ""}, fmt.Errorf("Couldn't tokenize Char from: %s", string(input[0]))
	}
}

func TokenizeKeywords(input string) (Token, error, int) {
	result := ""
	for i := 0; i < len(input); i++ {
		result += string(input[i])
		switch result {
		case "func":
			return Token{Type: 'F', Value: result}, nil, len(result) - 1
		case "i32":
			return Token{Type: 'T', Value: result}, nil, len(result) - 1
		case "i64":
			return Token{Type: 'T', Value: result}, nil, len(result) - 1
		case "f32":
			return Token{Type: 'T', Value: result}, nil, len(result) - 1
		case "void":
			return Token{Type: 'T', Value: result}, nil, len(result) - 1
		case "->":
			return Token{Type: 'r', Value: result}, nil, len(result) - 1
		case "return":
			return Token{Type: 'R', Value: result}, nil, len(result) - 1
		case "==":
			return Token{Type: 'E', Value: result}, nil, len(result) - 1
		}
	}
	return Token{Type: 'u', Value: ""}, fmt.Errorf("Couldn't tokenize any Keywords from: %s", string(input[0])), -1
}

func TokenizeName(input string) (Token, error, int) {
	pattern := regexp.MustCompile("[a-zA-Z]")
	result, err := TokenizePattern(pattern, input)
	if err != nil {
		panic("Something bad happened when tokenizing Names")
	}
	if result == "" {
		return Token{Type: 'u', Value: ""}, fmt.Errorf("Couldn't parse Name from: %s", input), -1
	}
	return Token{Type: 'N', Value: result}, nil, len(result) - 1
}

func Tokenize(input string) ([]Token, error) {
	var tokens []Token
	for i := 0; i < len(input); i++ {
		skip := SkipWhiteSpace(input[i:])
		i += skip
		token, err, n := TokenizeNumber(input[i:])
		if err == nil {
			tokens = append(tokens, token)
			i += n
			continue
		}
		token, err, n = TokenizeString(input[i:])
		if err == nil {
			tokens = append(tokens, token)
			i += n
			continue
		}
		token, err, n = TokenizeKeywords(input[i:])
		if err == nil {
			tokens = append(tokens, token)
			i += n
			continue
		}
		token, err, n = TokenizeName(input[i:])
		if err == nil {
			tokens = append(tokens, token)
			i += n
			continue
		}
		token, err = TokenizeChar(string(input[i]))
		if err == nil {
			tokens = append(tokens, token)
			continue
		}
		fmt.Printf("unrecognized char %s\n", string(input[i]))
	}
	return tokens, nil
}
