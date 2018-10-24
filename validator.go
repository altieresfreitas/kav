package main

import (
	"fmt"
	"net"
)

// Operators and Rules
const (
	MaskSize       RuleName = "MaskSize"
	ListSize       RuleName = "ListSize"
	LabelValues    RuleName = "LabelValues"
	OpIn           Operator = "In"
	OpNotIn        Operator = "NotIn"
	OpExists       Operator = "Exists"
	OpDoesNotExist Operator = "DoesNotExist"
	OpEq           Operator = "Equals"
	OpGt           Operator = "Gt"
	OpLt           Operator = "Lt"
)

// AdmissionValidator is a NetworkPolicy abstraction to validate objects
type AdmissionValidator struct {
	NetworkValidator NetworkValidator `json:"networkValidator"`
}

// NetworkValidator is ...
type NetworkValidator struct {
	PodSelector       PodSelector       `json:"podSelector"`
	NamespaceSelector NamespaceSelector `json:"namespaceSelector"`
	IPBlock           IPBlock           `json:"ipBlock"`
}

// PodSelector ..
type PodSelector struct {
	MatchLabels MatchLabels `json:"matchLabels"`
}

// NamespaceSelector ..
type NamespaceSelector struct {
	MatchLabels MatchLabels `json:"matchLabels"`
}

// MatchLabels ..
type MatchLabels struct {
	Rules []Rule `json:"rules"`
}

// IPBlock is ...
type IPBlock struct {
	Cidr   Cidr   `json:"cidr"`
	Except Except `json:"except"`
}

// Except is ..
type Except struct {
	Rules []Rule `json:"rules"`
}

// Cidr is ..
type Cidr struct {
	Rules []Rule `json:"rules"`
}

// RuleName is ..
type RuleName string

// Rule is ...
type Rule struct {
	Name     RuleName `json:"name"`
	Operator Operator `json:"operator"`
	Value    string   `json:"value"`
}

// validate ...
func (cidr *Cidr) validate() (bool, string, error) {
	var msg string

	_, _, err := net.ParseCIDR("")
	if err != nil {

		return false, msg, err
	}

	return false, "", nil
}

func check(x, y interface{}, o Operator) (bool, error) {

	switch o {
	case OpIn:
		fmt.Println("a")
	case OpNotIn:
		fmt.Println("b")
	case OpExists:
		fmt.Println("c")
	case OpDoesNotExist:
		fmt.Println("d")
	case OpGt:
		return checkOpGt(x.(int), y.(int))
	case OpLt:
		return checkOpLt(x.(int), y.(int))

	}

	return false, nil
}

// Operator is ..
type Operator string

func checkOpEq(x, y int) (bool, error) {

	if x == y {
		return true, nil
	}
	return false, nil
}

func checkOpGt(x, y int) (bool, error) {

	if x > y {
		return true, nil
	}
	return false, nil
}

func checkOpLt(x, y int) (bool, error) {

	if x < y {
		return true, nil
	}
	return false, nil
}
