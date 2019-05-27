package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
)

var stack = make([]int,1000)
var sp int

func push(val int) {
	stack[sp]=val
	sp++
}

func pop() (int, error) {
	if sp > 0 {
		sp--
		val := stack[sp]
		return val, nil
	} else {
		return 0, fmt.Errorf("Невозможно извлечь элемент из пустого стека")
	}
}

func empty() bool {
	return (sp == 0)
}

func calc(input io.Reader) (int, error) {
	in := bufio.NewScanner(input)
	in.Split(bufio.ScanWords)
	for in.Scan() {
		c := in.Text()
		switch c {
		case "\n": fallthrough
		case "=":
			if val, err := pop(); err == nil {
				return val, nil
			} else {
				return 0, err
			}
		case "+":
			if a, err := pop(); err == nil {
				if b, er := pop(); er == nil {
					push(a + b)
				} else {
					return 0, er
				}
			} else {
				return 0, err
			}
		case "-":
			if a, err := pop(); err == nil {
				if b, er := pop(); er == nil {
					push(b - a)
				} else {
					return 0, er
				}
			} else {
				return 0, err
			}
		case "*":
			if a, err := pop(); err == nil {
				if b, er := pop(); er == nil {
					push(a * b)
				} else {
					return 0, er
				}
			} else {
				return 0, err
			}
		case "/":
			if a, err := pop(); err == nil {
				if b, er := pop(); er == nil {
					push(b / a)
				} else {
					return 0, er
				}
			} else {
				return 0, err
			}
		default:
			if val, err := strconv.Atoi(c); err == nil {
				push(val)
			} else {
				return 0, err
			}
		}
	}
	return 0, fmt.Errorf("Конец файла")
}

func main() {
	if res, err := calc(os.Stdin); err == nil {
		fmt.Printf("Result = %d\n", res)
	} else {
		fmt.Println(err)
	}
	i := 0
	for !empty() {
		a, _ := pop()
		fmt.Printf("Stack[%d] = %d\n", i, a)
		i++
	}
}