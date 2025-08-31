package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Node interface {
	Evaluate() (float64, error)
	ToInfix() string
	ToPrefix() string
	ToPostfix() string
}

type NumberNode struct {
	Value float64
}

func (n *NumberNode) Evaluate() (float64, error) {
	return n.Value, nil
}

func (n *NumberNode) ToInfix() string   { return fmt.Sprintf("%.2f", n.Value) }
func (n *NumberNode) ToPrefix() string  { return fmt.Sprintf("%.2f", n.Value) }
func (n *NumberNode) ToPostfix() string { return fmt.Sprintf("%.2f", n.Value) }

type OperatorNode struct {
	Operator string
	Left     Node
	Right    Node
}

func (o *OperatorNode) Evaluate() (float64, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Critical error:", r)
		}
	}()

	leftVal, err := o.Left.Evaluate()
	if err != nil {
		return 0, err
	}
	rightVal, err := o.Right.Evaluate()
	if err != nil {
		return 0, err
	}

	switch o.Operator {
	case "+":
		return leftVal + rightVal, nil
	case "-":
		return leftVal - rightVal, nil
	case "*":
		return leftVal * rightVal, nil
	case "/":
		if rightVal == 0 {
			return 0, errors.New("division by zero")
		}
		return leftVal / rightVal, nil
	default:
		return 0, errors.New("unknown operator")
	}
}

func (o *OperatorNode) ToInfix() string {
	return fmt.Sprintf("(%s %s %s)", o.Left.ToInfix(), o.Operator, o.Right.ToInfix())
}

func (o *OperatorNode) ToPrefix() string {
	return fmt.Sprintf("%s %s %s", o.Operator, o.Left.ToPrefix(), o.Right.ToPrefix())
}

func (o *OperatorNode) ToPostfix() string {
	return fmt.Sprintf("%s %s %s", o.Left.ToPostfix(), o.Right.ToPostfix(), o.Operator)
}

type ParseError struct {
	Message string
}

func (e *ParseError) Error() string {
	return "Parse Error: " + e.Message
}

type ExpressionTree struct {
	Root Node
}

func ParseExpression(expression string) (*ExpressionTree, error) {
	var outputStack []Node
	var operatorStack []string
	tokens := strings.Fields(expression)

	precedence := map[string]int{
		"+": 1, "-": 1,
		"*": 2, "/": 2,
	}

	for _, token := range tokens {
		if value, err := strconv.ParseFloat(token, 64); err == nil {
			outputStack = append(outputStack, &NumberNode{Value: value})
		} else if token == "(" {
			operatorStack = append(operatorStack, token)
		} else if token == ")" {
			for len(operatorStack) > 0 && operatorStack[len(operatorStack)-1] != "(" {
				op := operatorStack[len(operatorStack)-1]
				operatorStack = operatorStack[:len(operatorStack)-1]
				if err := applyOperator(&outputStack, op); err != nil {
					return nil, err
				}
			}
			if len(operatorStack) == 0 {
				return nil, &ParseError{"Mismatched parentheses"}
			}
			operatorStack = operatorStack[:len(operatorStack)-1]
		} else {
			for len(operatorStack) > 0 && precedence[operatorStack[len(operatorStack)-1]] >= precedence[token] {
				op := operatorStack[len(operatorStack)-1]
				operatorStack = operatorStack[:len(operatorStack)-1]
				if err := applyOperator(&outputStack, op); err != nil {
					return nil, err
				}
			}
			operatorStack = append(operatorStack, token)
		}
	}

	for len(operatorStack) > 0 {
		op := operatorStack[len(operatorStack)-1]
		operatorStack = operatorStack[:len(operatorStack)-1]
		if err := applyOperator(&outputStack, op); err != nil {
			return nil, err
		}
	}

	if len(outputStack) != 1 {
		return nil, &ParseError{"Invalid expression"}
	}
	return &ExpressionTree{Root: outputStack[0]}, nil
}

func applyOperator(outputStack *[]Node, operator string) error {
	if len(*outputStack) < 2 {
		return &ParseError{"Not enough operands for operator " + operator}
	}
	right := (*outputStack)[len(*outputStack)-1]
	left := (*outputStack)[len(*outputStack)-2]
	*outputStack = (*outputStack)[:len(*outputStack)-2]
	*outputStack = append(*outputStack, &OperatorNode{Operator: operator, Left: left, Right: right})
	return nil
}

func (t *ExpressionTree) Print() {
	fmt.Println("Infix:", t.Root.ToInfix())
	fmt.Println("Prefix:", t.Root.ToPrefix())
	fmt.Println("Postfix:", t.Root.ToPostfix())
}

func (t *ExpressionTree) Evaluate() (float64, error) {
	return t.Root.Evaluate()
}

func main() {
	fmt.Println("Введіть арифметичний вираз (наприклад, 3 + 5 * ( 2 - 8 )):")
	var input string
	in := bufio.NewReader(os.Stdin)
	input, err := in.ReadString('\n')
	fmt.Println(input)

	tree, err := ParseExpression(input)
	if err != nil {
		fmt.Println("Помилка парсингу виразу:", err)
		return
	}

	tree.Print()
	result, err := tree.Evaluate()
	if err != nil {
		fmt.Println("Помилка обчислення:", err)
	} else {
		fmt.Println("Результат обчислення:", result)
	}
}
