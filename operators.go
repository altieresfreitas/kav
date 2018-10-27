package main

import "fmt"

// Operator is validation type abstraction
type Operator string

// All operators available
const (
	OpIn           Operator = "In"
	OpNotIn        Operator = "NotIn"
	OpExists       Operator = "Exists"
	OpDoesNotExist Operator = "DoesNotExist"
	OpEq           Operator = "Equals"
	OpGt           Operator = "Gt"
	OpLt           Operator = "Lt"
	OpGe           Operator = "Ge"
	OpLe           Operator = "Le"
)

func operatorExec(x, y interface{}, o Operator) (bool, error) {

	switch o {

	case OpIn:

		fmt.Println("a")

	case OpNotIn:

		fmt.Println("b")

	case OpExists:

		fmt.Println("c")

	case OpDoesNotExist:

		fmt.Println("d")

	case OpGe:

		return opGeExec(x.(int), y.(int)), nil

	case OpGt:

		return opGtExec(x.(int), y.(int))

	case OpLt:

		return opLtExec(x.(int), y.(int))

	case OpLe:

		return opLeExec(x.(int), y.(int)), nil

	case OpEq:

		return opEqExec(x.(int), y.(int)), nil

	}

	return false, nil
}

func opEqExec(x, y int) bool {

	if x == y {
		return true
	}
	return false
}

func opGtExec(x, y int) (bool, error) {

	if x > y {
		return true, nil
	}
	return false, nil
}

func opLtExec(x, y int) (bool, error) {

	if x < y {
		return true, nil
	}
	return false, nil
}

func opGeExec(x, y int) bool {

	if x >= y {
		return true
	}
	return false
}

func opLeExec(x, y int) bool {

	if x <= y {
		return true
	}
	return false
}
