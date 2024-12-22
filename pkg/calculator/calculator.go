package calculator

import (
	"errors"
	"strconv"
	"unicode"
)

var (
	ErrDivisionByZero    = errors.New("division by zero")
	ErrUnknownOperator   = errors.New("unknown operator")
	ErrExpressionInvalid = errors.New("expression is invalid")
)

var OPERATION_PRIORITIES = map[rune]int{
	'*': 1,
	'/': 1,
	'-': 2,
	'+': 2,
	'(': 3,
}

func Pop[T any](array *[]T) (T, error) {
	if len(*array) == 0 {
		var zeroVar T
		return zeroVar, ErrExpressionInvalid
	}

	elem := (*array)[len(*array)-1]
	*array = (*array)[:len(*array)-1]

	return elem, nil
}

func compute(a, b float64, operator string) (float64, error) {
	switch operator {
	case "+":
		return a + b, nil
	case "-":
		return a - b, nil
	case "*":
		return a * b, nil
	case "/":
		if b == 0 {
			return 0, ErrDivisionByZero
		}
		return a / b, nil
	default:
		return 0, ErrUnknownOperator
	}
}

func checkExpression(expression string) bool {
	for _, symbol := range expression {
		_, check := OPERATION_PRIORITIES[symbol]
		if !unicode.IsDigit(symbol) &&
			!check &&
			symbol != ' ' &&
			symbol != '(' &&
			symbol != ')' {
			return false
		}
	}
	return true
}

func toPostfix(expression string) ([]string, error) {
	stack := []rune{}
	output := []string{}
	expressionLength := len(expression)
	index := 0
	operSign := "+"

	isUnary := true
	before_bracket := false

	for index < expressionLength {
		symbol := rune(expression[index])

		if symbol != ' ' {
			if unicode.IsDigit(symbol) {
				cur := index
				for cur < expressionLength-1 && unicode.IsDigit(rune(expression[cur+1])) {
					cur++
				}

				add := ""
				if operSign == "-" {
					add = "-"
					operSign = "+"
				}

				if cur == expressionLength-1 {
					output = append(output, add+expression[index:])
				} else {
					output = append(output, add+expression[index:cur+1])
				}
				index = cur
				isUnary = false
			} else if symbol == '(' {
				if operSign == "-" {
					operSign = "+"
					before_bracket = true
				}
				stack = append(stack, symbol)
				isUnary = true
			} else if symbol == ')' {
				for len(stack) > 0 && stack[len(stack)-1] != '(' {
					op, err := Pop(&stack)
					if err != nil {
						return []string{}, err
					}
					output = append(output, string(op))
				}
				if before_bracket {
					output = append(output, "-")
				}
				_, err := Pop(&stack)
				if err != nil {
					return []string{}, err
				}
				isUnary = false
			} else {
				if symbol == '-' || symbol == '+' {
					if isUnary {
						operSign = string(symbol)
					} else {
						for len(stack) > 0 && OPERATION_PRIORITIES[stack[len(stack)-1]] <= OPERATION_PRIORITIES[symbol] {
							op, err := Pop(&stack)
							if err != nil {
								return []string{}, err
							}
							output = append(output, string(op))
						}
						stack = append(stack, symbol)
					}
					index++
					isUnary = true
					continue
				}

				for len(stack) > 0 && OPERATION_PRIORITIES[stack[len(stack)-1]] <= OPERATION_PRIORITIES[symbol] {
					op, err := Pop(&stack)
					if err != nil {
						return []string{}, err
					}
					output = append(output, string(op))
				}
				stack = append(stack, symbol)
				isUnary = true
			}
		}
		index++
	}

	for len(stack) > 0 {
		op, _ := Pop(&stack)
		output = append(output, string(op))
	}
	return output, nil
}

func Calc(expression string) (float64, error) {
	if !checkExpression(expression) {
		return 0, ErrExpressionInvalid
	}
	postfixExpression, err := toPostfix(expression)
	if err != nil {
		return 0, err
	}
	stack := []float64{}

	for _, item := range postfixExpression {
		isOperator := false
		for _, op := range []string{"/", "*", "-", "+"} {
			if item == op {
				isOperator = true
				break
			}
		}
		if isOperator {
			b, err := Pop(&stack)
			if err != nil {
				return 0, err
			}
			a, err := Pop(&stack)
			if err != nil {
				if item == "+" || item == "-" {
					a = 0
				} else {
					return 0, err
				}
			}

			result, err := compute(a, b, item)
			if err != nil {
				return 0, err
			}
			stack = append(stack, result)
		} else {
			value, err := strconv.ParseFloat(item, 64)
			if err != nil {
				return 0, err
			}

			stack = append(stack, value)
		}
	}

	if len(stack) == 1 {
		return stack[0], nil
	}
	return 0, ErrExpressionInvalid
}
