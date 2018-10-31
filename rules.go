package main

import (
	"fmt"
	"net"

	"k8s.io/apimachinery/pkg/util/intstr"
)

// RuleName is ..
type RuleName string

// Rules
const (
	MaskBitsSize RuleName = "MaskBitsSize"
	ListSize     RuleName = "ListSize"
	LabelValues  RuleName = "LabelValues"
	LabelCount   RuleName = "LabelCount"
	PortNumber   RuleName = "PortNumber"
)

// Rule is ...
type Rule struct {
	Name     RuleName           `json:"name"`
	Operator Operator           `json:"operator"`
	Key      string             `json:"key"`
	Value    intstr.IntOrString `json:"value"`
}

func (v *Rule) isValidMask(c string) (bool, error) {

	_, network, _ := net.ParseCIDR(c)
	m, _ := network.Mask.Size()

	ok, err := operatorExec(m, v.Value.IntValue(), v.Operator)
	if err != nil {
		return false, err
	}

	if !ok {
		return false, fmt.Errorf(
			"error InvalidMaskSize: mask size must be %s %v",
			v.Operator, v.Value.IntValue())
	}

	return true, nil

}

func (v *Rule) isValidMaskList(l []string) (bool, error) {

	for _, i := range l {

		ok, err := v.isValidMask(i)
		if err != nil {
			return false, err
		}

		if !ok {
			return false, fmt.Errorf(
				"error InvalidMaskSize: mask size must be %s %v",
				v.Operator, v.Value.IntValue())
		}

	}

	return true, nil

}

func (v *Rule) isValidListSize(s int) (bool, error) {

	ok, err := operatorExec(s, v.Value.IntValue(), v.Operator)
	if err != nil {
		return false, err
	}

	if !ok {
		return false, fmt.Errorf(
			"error InvalidListSize: list size must be %s %v ",
			v.Operator, v.Value.IntValue())
	}

	return true, nil

}

func (v *Rule) isValidPort(s int) (bool, error) {

	ok, err := operatorExec(s, v.Value.IntValue(), v.Operator)
	if err != nil {
		return false, err
	}

	if !ok {
		return false, fmt.Errorf(
			"error InvalidPortNumber: port number must be %s %v ",
			v.Operator, v.Value.IntValue())
	}

	return true, nil

}

func (v *Rule) isValidLabelValues(labels map[string]string) (bool, error) {

	ok, err := operatorExec(len(labels), v.Value.IntValue(), v.Operator)
	if err != nil {
		return false, err
	}

	if !ok {
		return false, fmt.Errorf(
			"error InvalidLabelCount: the numbers of labels must be %s %v",
			v.Operator, v.Value.IntValue())
	}

	return true, nil

}

func (v *Rule) isValidLabelCount(labels map[string]string) (bool, error) {

	ok, err := operatorExec(len(labels), v.Value.IntValue(), v.Operator)
	if err != nil {
		return false, err
	}

	if !ok {
		return false, fmt.Errorf(
			"error InvalidLabelCount: the numbers of labels must be %s %v",
			v.Operator, v.Value.IntValue())
	}

	return true, nil

}
