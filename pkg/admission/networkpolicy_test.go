package admission

import (
	"encoding/json"
	"testing"

	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestCIDR(t *testing.T) {

	tests := []struct {
		rules    string
		expected bool
		CIDR     string
	}{
		{
			`{ "rules": [
				{
					"name": "MaskBitsSize",
					"operator": "Gt",
					"value": 20
				}
			]}`,
			true,
			`{ "cidr":"192.168.0.0/24"}`,
		},
		{
			`{ "rules": [
				{
					"name": "MaskBitsSize",
					"operator": "Gt",
					"value": 20
				}
			]}`,
			false,
			`{ "cidr":"192.168.0.0/19"}`,
		},
		{
			`{ "rules": [
				{
					"name": "MaskBitsSize",
					"operator": "Lt",
					"value": 20
				}
			]}`,
			true,
			`{ "cidr":"192.168.0.0/19"}`,
		},
		{
			`{ "rules": [
				{
					"name": "MaskBitsSize",
					"operator": "Equals",
					"value": 20
				}
			]}`,
			true,
			`{ "cidr":"192.168.0.0/20"}`,
		},
	}

	for _, i := range tests {
		a := CIDR{}
		b := networkingv1.IPBlock{}

		if err := json.Unmarshal([]byte(i.rules), &a); err != nil {
			t.Errorf("error %v", err)
		}

		if err := json.Unmarshal([]byte(i.CIDR), &b); err != nil {
			t.Errorf("error %v", err)
		}

		result, _ := a.isValid(&b)

		if result != i.expected {
			t.Errorf(" %v", result)
		}

	}

}

func TestAllowedPolicyType(t *testing.T) {

	tests := []struct {
		allowedPolicyTypes string
		expected           bool
		policyTypes        string
	}{
		{
			`{"allowedPolicyTypes": [
				"Ingress",
				"Egress"]
				}`,
			true,
			`{
				"spec": {
				  "podSelector": {
					"matchLabels": {
					  "role": "db"
					}
				  },
				  "policyTypes": [
					"Ingress",
					"Egress"
				  ],
				  "ingress": [
					{
					  "from": [
						{
						  "ipBlock": {
							"cidr": "172.17.0.0/29"
						  }
						}
					  ]
					}
				  ]
				}
			  }`,
		},
		{
			`{"allowedPolicyTypes": [ "Egress" ] }`,
			false,
			`{
				"spec": {
					"policyTypes": [
						"Ingress",
						"Egress"
					  ],
				  "ingress": [
					{
					  "from": [
						{
						  "ipBlock": {
							"cidr": "172.17.0.0/29"
						  }
						}
					  ]
					}
				  ]
				}
			  }`,
		},
	}

	for _, i := range tests {
		a := NetworkPolicyValidator{}
		b := networkingv1.NetworkPolicy{}

		if err := json.Unmarshal([]byte(i.allowedPolicyTypes), &a); err != nil {
			t.Errorf("error %v", err)
		}

		if err := json.Unmarshal([]byte(i.policyTypes), &b); err != nil {
			t.Errorf("error %v", err)
		}

		result, _ := a.isValid(&b.Spec)

		if result != i.expected {
			t.Errorf(" %v", result)
		}

	}

}

func TestExcept(t *testing.T) {

	tests := []struct {
		rules    string
		expected bool
		except   string
	}{
		{
			`{ "rules": [
				{
					"name": "ListSize",
					"operator": "Ge",
					"value": 1
				}
			]}`,
			true,
			`{ "except":["192.168.0.0/19"]}`,
		},
		{
			`{ "rules": [
				{
					"name": "ListSize",
					"operator": "Le",
					"value": 1
				}
			]}`,
			true,
			`{ "except":["192.168.0.0/19"]}`,
		},
		{
			`{ "rules": [
				{
					"name": "ListSize",
					"operator": "Lt",
					"value": 1
				}
			]}`,
			false,
			`{ "except":["192.168.0.0/19"]}`,
		},
		{
			`{ "rules": [
				{
					"name": "ListSize",
					"operator": "Equals",
					"value": 2
				}
			]}`,
			true,
			`{ "except":["192.168.0.0/19","192.168.0.0/19"]}`,
		},
	}

	for _, i := range tests {
		a := Except{}
		b := networkingv1.IPBlock{}

		if err := json.Unmarshal([]byte(i.rules), &a); err != nil {
			t.Errorf("error %v", err)
		}

		if err := json.Unmarshal([]byte(i.except), &b); err != nil {
			t.Errorf("error %v", err)
		}

		result, err := a.isValid(&b)

		if result != i.expected {
			t.Errorf(" %v", err)
		}

	}

}

func TestMatchLabels(t *testing.T) {

	tests := []struct {
		rules       string
		expected    bool
		MatchLabels string
	}{
		{`{ "rules": [
			{
				"name": "LabelCount",
				"operator": "Equals",
				"value": 1
			}
		]}`,
			true,
			`{
				"matchLabels": {
				"teste": "db"
				}
			}`,
		},
		{`{ "rules": [
			{
				"name": "LabelCount",
				"operator": "Ge",
				"value": 1
			}
		]}`,
			true,
			`{
				"matchLabels": {
				"teste": "db"
				}
			}`,
		},
		{`{ "rules": [
			{
				"name": "LabelCount",
				"operator": "Le",
				"value": 1
			}
		]}`,
			true,
			`{
				"matchLabels": {
				"teste": "db"
				}
			}`,
		},
		{`{ "rules": [
			{
				"name": "LabelCount",
				"operator": "Lt",
				"value": 1
			}
		]}`,
			false,
			`{
				"matchLabels": {
				"teste": "db"
				}
			}`,
		},
		{`{"rules": [
			{
				"name": "LabelCount",
				"operator": "Gt",
				"value": 3
			},
			{
				"name": "LabelValues",
				"operator": "DoesNotExist",
				"value": ""
			}
		]}`,
			false,
			`{
				"matchLabels": {
				"teste": "db"
				}
			}`,
		},
		{`{"rules": [
			{
				"name": "LabelCount",
				"operator": "Gt",
				"value": 3
			},
			{
				"name": "LabelValues",
				"operator": "DoesNotExist",
				"value": ""
			}
		]}`,
			false,
			`{
				"matchLabels": {
				"teste": "db"
				}
			}`,
		},
	}

	for _, i := range tests {
		a := MatchLabels{}
		b := metav1.LabelSelector{}

		if err := json.Unmarshal([]byte(i.rules), &a); err != nil {
			t.Errorf("error %v", err)
		}

		if err := json.Unmarshal([]byte(i.MatchLabels), &b); err != nil {
			t.Errorf("error %v", err)
		}

		result, _ := a.isValid(b.MatchLabels)

		if result != i.expected {
			t.Errorf(" %v", result)
		}

	}

}
