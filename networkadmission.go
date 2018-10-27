package main

import (
	"net"
	"regexp"

	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NetworkAdmissionValidator is a NetworkPolicy abstraction to isValid objects
type NetworkAdmissionValidator struct {
	NetworkValidator []NetworkPolicyValidator `json:"networkValidator,omitempty"`
}

// IsValid will compare a received network policy object with NetworkadmissionRules.
func (v *NetworkAdmissionValidator) IsValid(p *networkingv1.NetworkPolicy) (bool, error) {

	for _, i := range v.NetworkValidator {

		r, err := regexp.Compile(i.Namespace)
		if err != nil {

			return false, err

		}

		if r.Match([]byte(p.Namespace)) {

			if ok, err := i.isValid(&p.Spec); !ok {

				return false, err

			}
		}
	}
	return true, nil

}

// PolicyType is a List of allowed Policies type, e.g Ingress or Egress
type PolicyType networkingv1.PolicyType

const (
	// PolicyTypeIngress is a NetworkPolicy that affects ingress traffic on selected pods
	PolicyTypeIngress PolicyType = "Ingress"
	// PolicyTypeEgress is a NetworkPolicy that affects egress traffic on selected pods
	PolicyTypeEgress PolicyType = "Egress"
)

// NetworkPolicyValidator provides the specification of a NetworkPolicy
type NetworkPolicyValidator struct {
	Namespace   string                   `json:"namespace,omitempty"`
	Ingress     NetworkPolicyIngressRule `json:"ingress,omitempty"`
	Egress      NetworkPolicyEgressRule  `json:"egress,omitempty"`
	PolicyTypes []PolicyType             `json:"policyTypes,omitempty"`
	PodSelector PodSelector              `json:"podSelector,omitempty"`
}

func (v *NetworkPolicyValidator) isValid(p *networkingv1.NetworkPolicySpec) (bool, error) {

	/* 	if ok, err := v.isValidPolicyTypes(&p.PolicyTypes); !ok {
	   		return false, err
	   	}
	*/
	if ok, err := v.PodSelector.isValid(&p.PodSelector); !ok {
		return false, err
	}

	if ok, err := v.Egress.isValid(&p.Egress); !ok {
		return false, err
	}

	if ok, err := v.Ingress.isValid(&p.Ingress); !ok {
		return false, err
	}

	return true, nil

}

/* func (v *NetworkPolicyValidator) isValidPolicyTypes(n *[]networkingv1.PolicyType) (bool, error) {

	for _, i := range *n {
		for _, j := range v.PolicyTypes {
			fmt.Println(i)
			fmt.Println(j)
		}

	}
	return false, nil
} */

// NetworkPolicyIngressRule describes a particular set of traffic that is allowed to the pods
type NetworkPolicyIngressRule struct {
	Ports NetworkPolicyPort `json:"ports,omitempty"`
	From  NetworkPolicyPeer `json:"from,omitempty"`
}

func (v *NetworkPolicyIngressRule) isValid(p *[]networkingv1.NetworkPolicyIngressRule) (bool, error) {

	for _, e := range *p {

		if ok, err := v.From.isValid(e.From); !ok {
			return false, err
		}

	}

	return true, nil
}

// NetworkPolicyEgressRule describes a particular set of traffic that is allowed out of pods
type NetworkPolicyEgressRule struct {
	Ports []NetworkPolicyPort `json:"ports,omitempty"`
	To    NetworkPolicyPeer   `json:"to,omitempty"`
}

func (v *NetworkPolicyEgressRule) isValid(p *[]networkingv1.NetworkPolicyEgressRule) (bool, error) {

	for _, e := range *p {

		if ok, err := v.To.isValid(e.To); !ok {
			return false, err
		}

	}

	return true, nil
}

// NetworkPolicyPort describes a port to allow traffic on
type NetworkPolicyPort struct {
	Protocol Protocol `json:"protocol,omitempty"`
	Port     Port     `json:"port,omitempty"`
}

func (p *NetworkPolicyPort) isValid(er *networkingv1.NetworkPolicyPort) (bool, error) {

	return true, nil

}

// Protocol ssad
type Protocol struct {
	Rules []Rule `json:"rules"`
}

// Port is
type Port struct {
	Rules []Rule `json:"rules"`
}

// IPBlock describes a particular CIDR (Ex. "192.168.1.1/24") that is allowed to the pods
type IPBlock struct {
	CIDR   CIDR   `json:"cidr"`
	Except Except `json:"except,omitempty"`
}

func (v *IPBlock) isValid(p *networkingv1.IPBlock) (bool, error) {

	if ok, err := v.CIDR.isValid(p); !ok {
		return false, err
	}

	if ok, err := v.Except.isValid(p); !ok {
		return false, err
	}

	return true, nil
}

// CIDR abstract CIDR object
type CIDR struct {
	Rules []Rule `json:"rules"`
}

func (c *CIDR) isValid(p *networkingv1.IPBlock) (bool, error) {

	_, network, _ := net.ParseCIDR(p.CIDR)
	maskBits, _ := network.Mask.Size()

	for _, r := range c.Rules {

		switch r.Name {

		case MaskBitsSize:
			return r.isValidMaskBitsSize(maskBits)

		}
	}

	return true, nil

}

// Except abstract CIDR object
type Except struct {
	Rules []Rule `json:"rules"`
}

func (v *Except) isValid(p *networkingv1.IPBlock) (bool, error) {

	for _, r := range v.Rules {

		switch r.Name {

		case ListSize:

			return r.isValidListSize(p.Except)

		case MaskBitsSize:

			for _, n := range p.Except {

				_, network, _ := net.ParseCIDR(n)
				maskBits, _ := network.Mask.Size()

				if ok, err := r.isValid(maskBits); !ok {

					return false, err

				}

			}

		}
	}

	return true, nil
}

// NetworkPolicyPeer describes a peer to allow traffic from. Only certain combinations of
type NetworkPolicyPeer struct {
	PodSelector       PodSelector       `json:"podSelector,omitempty"`
	NamespaceSelector NamespaceSelector `json:"namespaceSelector,omitempty"`
	IPBlock           IPBlock           `json:"ipBlock,omitempty"`
}

func (v *NetworkPolicyPeer) isValid(p []networkingv1.NetworkPolicyPeer) (bool, error) {

	for _, i := range p {

		switch {

		case i.PodSelector != nil:

			if ok, err := v.PodSelector.isValid(i.PodSelector); !ok {
				return false, err
			}

		case i.NamespaceSelector != nil:

			if ok, err := v.NamespaceSelector.isValid(i.NamespaceSelector); !ok {
				return false, err
			}

		case i.IPBlock != nil:

			if ok, err := v.IPBlock.isValid(i.IPBlock); !ok {
				return false, err
			}

		}

	}

	return true, nil

}

// PodSelector ..
type PodSelector struct {
	MatchLabels MatchLabels `json:"matchLabels"`
}

func (v *PodSelector) isValid(p *metav1.LabelSelector) (bool, error) {

	return v.MatchLabels.isValid(p.MatchLabels)

}

// NamespaceSelector ..
type NamespaceSelector struct {
	MatchLabels MatchLabels `json:"matchLabels"`
}

func (v *NamespaceSelector) isValid(p *metav1.LabelSelector) (bool, error) {

	return v.MatchLabels.isValid(p.MatchLabels)

}

// MatchLabels ..
type MatchLabels struct {
	Rules []Rule `json:"rules"`
}

func (v *MatchLabels) isValid(p map[string]string) (bool, error) {

	for _, r := range v.Rules {

		switch r.Name {

		case LabelCount:
			return r.isValidLabelCount(p)

		}
	}

	return true, nil

}
