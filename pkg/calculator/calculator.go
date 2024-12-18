package calculator

import (
	"errors"
	"fmt"
	"math/bits"
	"strconv"
	"strings"
)

// Calc evaluates the value of an arithmetic expression.
func Calc(expression string) (float64, error) {
	tokens := tokenize(expression)
	if len(tokens) == 0 {
		return 0, errors.New("the expression is empty")
	}

	var numStack []float64
	var opStack []string
	var err error
	for _, token := range tokens {
		if isNumber(token) {
			num, err := strconv.ParseFloat(token, bits.UintSize)
			if err != nil {
				return 0, fmt.Errorf("invalid number: '%s'", token)
			}
			numStack = append(numStack, num)
		} else if token == "(" {
			opStack = append(opStack, token)
		} else if token == ")" {
			for len(opStack) > 0 && opStack[len(opStack)-1] != "(" {
				numStack, opStack, err = evaluate(numStack, opStack)
				if err != nil {
					return 0, err
				}
			}
			if len(opStack) == 0 {
				return 0, errors.New("unmatched closing parenthesis")
			}
			opStack = opStack[:len(opStack)-1]
		} else if isOperator(token) {
			for len(opStack) > 0 && precedence(opStack[len(opStack)-1]) >= precedence(token) {
				numStack, opStack, err = evaluate(numStack, opStack)
				if err != nil {
					return 0, err
				}
			}
			opStack = append(opStack, token)
		} else {
			return 0, fmt.Errorf("unknown token: '%s'", token)
		}
	}

	for len(opStack) > 0 {
		numStack, opStack, err = evaluate(numStack, opStack)
		if err != nil {
			return 0, err
		}
	}

	if len(numStack) != 1 {
		return 0, errors.New("invalid expression: check the number of operators and operands")
	}

	return numStack[0], nil
}

// tokenize splits the expression into tokens.
func tokenize(expr string) []string {
	var tokens []string
	var currentToken strings.Builder

	for _, char := range expr {
		if char == ' ' {
			continue
		}
		if isOperator(string(char)) || char == '(' || char == ')' {
			if currentToken.Len() > 0 {
				tokens = append(tokens, currentToken.String())
				currentToken.Reset()
			}
			tokens = append(tokens, string(char))
		} else {
			currentToken.WriteRune(char)
		}
	}
	if currentToken.Len() > 0 {
		tokens = append(tokens, currentToken.String())
	}

	return tokens
}

func isNumber(s string) bool {
	_, err := strconv.ParseFloat(s, bits.UintSize)
	return err == nil
}

func isOperator(s string) bool {
	return s == "+" || s == "-" || s == "*" || s == "/"
}

func precedence(op string) int {
	switch op {
	case "+", "-":
		return 1
	case "*", "/":
		return 2
	}
	return 0
}

func evaluate(numStack []float64, opStack []string) ([]float64, []string, error) {
	if len(numStack) < 2 {
		return numStack, opStack, errors.New("not enough operands for the operation")
	}
	if len(opStack) == 0 {
		return numStack, opStack, errors.New("no operators available")
	}

	b := numStack[len(numStack)-1]
	a := numStack[len(numStack)-2]
	op := opStack[len(opStack)-1]

	numStack = numStack[:len(numStack)-2]
	opStack = opStack[:len(opStack)-1]

	var result float64
	switch op {
	case "+":
		result = a + b
	case "-":
		result = a - b
	case "*":
		result = a * b
	case "/":
		if b == 0 {
			return numStack, opStack, errors.New("division by zero is not allowed")
		}
		result = a / b
	default:
		return numStack, opStack, errors.New("unknown operator")
	}

	numStack = append(numStack, result)
	return numStack, opStack, nil
}
