package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

func getNetworkValidator() *NetworkAdmissionValidator {
	var policyFile = "" //flags.Configfile

	_, err := os.Stat(policyFile)
	if err != nil {
		log.Fatal("Config file is missing: ", policyFile)
	}

	yamlFile, err := os.Open(policyFile)
	if err != nil {
		fmt.Println(err)
	}

	defer yamlFile.Close()
	byteValue, _ := ioutil.ReadAll(yamlFile)
	jsonFile, err := yaml.ToJSON(byteValue)
	if err != nil {
		fmt.Println(err)
	}

	policy := &NetworkAdmissionValidator{}
	json.Unmarshal(jsonFile, &policy)

	return policy
}

func getNetworkPolicy() *networkingv1.NetworkPolicy {
	var policyFile = "" //flags.Configfile

	_, err := os.Stat(policyFile)
	if err != nil {
		log.Fatal("Config file is missing: ", policyFile)
	}

	yamlFile, err := os.Open(policyFile)
	if err != nil {
		fmt.Println(err)
	}

	defer yamlFile.Close()
	byteValue, _ := ioutil.ReadAll(yamlFile)
	jsonFile, err := yaml.ToJSON(byteValue)
	if err != nil {
		fmt.Println(err)
	}
	policy := &networkingv1.NetworkPolicy{}
	json.Unmarshal(jsonFile, &policy)

	return policy
}

func main() {

	a := getNetworkPolicy()
	b := getNetworkValidator()

	if ok, err := b.IsValid(a); !ok {

		fmt.Println(err)

	}

}
