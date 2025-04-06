package main

import "fmt"

type Employee struct {
	Name   string
	Salary int
}

// Function to give a raise to an employee
func giveRaise(e *Employee, raise int) {
	e.Salary += raise
}

func main() {
	emp := Employee{Name: "John Doe", Salary: 50000}

	giveRaise(&emp, 5000)
	fmt.Println("After raise:", emp)
}
