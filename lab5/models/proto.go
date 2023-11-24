package models

import (
	"fmt"
	"math"
)

// Message - структура, которая описывает то, как будут передаваться
// данные с клиента на сервер
type Message struct {
	A float64 `json:"a"`
	B float64 `json:"b"`
	C float64 `json:"c"`
}

// InputFromConsole - функция для ввода структуры Message с консоли
func (m Message) InputFromConsole() Message {
	fmt.Print("Enter a sequence of numbers (separated by spaces): ")
	a := 0.0
	b := 0.0
	c := 0.0

	fmt.Scanln(&a, &b, &c)
	return Message{
		A: a,
		B: b,
		C: c,
	}
}

// Solve - вспомогательная стуруктура для Response
type Solve struct {
	Root float64
}

// Response - структура, которая описывает, как данные передаются с сервера обратно на клиент
// то есть это просто структура ответа на запрос
type Response struct {
	Roots []Solve
}

// OutToConsole - функция для ввода структуры Message в консоль
func (r Response) OutToConsole() {
	fmt.Println("Решения")
	for _, v := range r.Roots {
		fmt.Print(v)
	}
	fmt.Println()
}

// ResultFunction - функция, которая возвращает результат, то есть
// это главная функция, занимающаяся вычислениями / выполнением каких-то действий
func ResultFunction(numbers Message) (Response, error) {
	discriminant := numbers.B*numbers.B - 4*numbers.A*numbers.C

	roots := make([]Solve, 0)

	if discriminant > 0 {
		root1 := (-numbers.B + math.Sqrt(discriminant)) / (2 * numbers.A)
		root2 := (-numbers.B - math.Sqrt(discriminant)) / (2 * numbers.A)

		roots = append(roots, Solve{Root: root1})
		roots = append(roots, Solve{Root: root2})
	} else if discriminant == 0 {
		root := -numbers.B / (2 * numbers.A)
		roots = append(roots, Solve{Root: root})
	}

	return Response{Roots: roots}, nil
}
