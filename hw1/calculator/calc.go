package main

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
	"unicode"
)

func precedence(op rune) int {
	switch op {
	case '+', '-':
		return 1
	case '*', '/':
		return 2
	}
	return 0
}

func isValidNumber(num string) bool {
	dotCount := 0
	for index, char := range num {
		if !unicode.IsDigit(char) {
			if (char == '-' || char == '+') && index == 0 {
				continue
			}
			if char == '.' {
				dotCount++
				if dotCount > 1 {
					return false // больше одной точки
				}
			} else {
				return false // некорректный символ
			}
		}
	}
	return num[len(num)-1] != '.'
}

func infixToPostfix(expression string) ([]string, error) {
	var result []string
	var stack []rune
	if len(expression) == 0 {
		return result, nil
	}

	var currentNum string
	lastWasOperator := true // Указывает, был ли последний символ оператором или скобкой
	if (expression[0] == '-' || expression[0] == '+') && (len(expression) > 1 && expression[1] == '(') {
		expression = "0" + expression
	}
	for _, char := range expression {
		if unicode.IsSpace(char) {
			continue // Игнорируем пробелы
		}

		if unicode.IsDigit(char) || char == '.' {
			currentNum += string(char) // Собираем число
			lastWasOperator = false    // Последний символ - не оператор
		} else {
			if currentNum != "" {
				if !isValidNumber(currentNum) {
					return nil, errors.New("Invalid number: " + currentNum)
				}
				result = append(result, currentNum)
				currentNum = "" // Сбрасываем текущее число
			}

			switch char {
			case '+', '-':
				if lastWasOperator { // Обработка унарного оператора
					if currentNum == "" {
						currentNum = string(char)
					} else {
						return nil, errors.New("Invalid input")
					}
					lastWasOperator = true
					continue
				}
				for len(stack) > 0 && precedence(stack[len(stack)-1]) >= precedence(char) {
					result = append(result, string(stack[len(stack)-1]))
					stack = stack[:len(stack)-1]
				}
				stack = append(stack, char)
				lastWasOperator = true // Устанавливаем, что последний символ - оператор
			case '*', '/':
				for len(stack) > 0 && precedence(stack[len(stack)-1]) >= precedence(char) {
					result = append(result, string(stack[len(stack)-1]))
					stack = stack[:len(stack)-1]
				}
				stack = append(stack, char)
				lastWasOperator = true // Устанавливаем, что последний символ - не оператор
			case '(':
				stack = append(stack, '(')
				lastWasOperator = true // '(' может быть перед унарным оператором
			case ')':
				for len(stack) > 0 && stack[len(stack)-1] != '(' {
					result = append(result, string(stack[len(stack)-1]))
					stack = stack[:len(stack)-1]
				}
				if len(stack) == 0 {
					return nil, errors.New("Mismatched parentheses")
				}
				stack = stack[:len(stack)-1] // убираем '(' (закрываем скобки)
				lastWasOperator = false
			default:
				return nil, errors.New("Invalid character: " + string(char))
			}
		}
	}

	if currentNum != "" {
		if !isValidNumber(currentNum) {
			return nil, errors.New("Invalid number: " + currentNum)
		}
		result = append(result, currentNum)
	}

	for len(stack) > 0 {
		result = append(result, string(stack[len(stack)-1]))
		stack = stack[:len(stack)-1]
	}
	return result, nil
}

func isOperator(char rune) bool {
	return char == '+' || char == '-' || char == '*' || char == '/'
}

func calcExpr(expression string) (float64, error) {
	postFixExp, err := infixToPostfix(expression)
	if err != nil {
		return 0, err
	}

	stack := []float64{}

	for _, token := range postFixExp {

		if isOperator(rune(token[len(token)-1])) {
			if len(stack) <= 1 {
				return 0, errors.New("Invalid expression")
			}
			secondNum := stack[len(stack)-1]
			firstNum := stack[len(stack)-2]
			stack = stack[:len(stack)-2] // Убираем последние два числа

			var result float64
			switch token {
			case "+":
				result = firstNum + secondNum
			case "-":
				result = firstNum - secondNum
			case "*":
				result = firstNum * secondNum
			case "/":
				if secondNum == 0 {
					return 0, errors.New("Division by zero")
				}
				result = firstNum / secondNum
			}
			stack = append(stack, result)
		} else {
			num, err := strconv.ParseFloat(token, 64)
			if err != nil {
				return 0, err
			}
			stack = append(stack, num)
		}
	}

	if len(stack) != 1 {
		return 0, errors.New("Invalid expression")
	}
	roundTo := 10000
	return math.Ceil(stack[0]*float64(roundTo)) / float64(roundTo), nil
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Введите строку: ")
	for scanner.Scan() {

		expression := scanner.Text()
		expAnswer, err := calcExpr(expression)
		if err == nil {
			fmt.Println(expAnswer)
		} else {
			fmt.Println(err)
		}
		fmt.Print("Введите строку: ")
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Ошибка при чтении:", err)
	}

}
