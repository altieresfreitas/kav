# karepol - Kubernetes Admission Resource Policy

This is a generic webhook admission controller to configure resource restrictions on top of Kubernetes.

Currently only networkpolicy resource is supported

## How It Works
1 - Enable the dynamic admission controller registration API by adding admissionregistration.k8s.io/v1alpha1 to the --runtime-config flag passed to kube-apiserver, e.g. --runtime-config=admissionregistration.k8s.io/v1alpha1. Again, all replicas should have the same flag setting.

2 - Create ValidatingWebhookConfiguration
kubectl apply -f examples/k8s-validator.yaml

```
apiVersion: admissionregistration.k8s.io/v1beta1
kind: ValidatingWebhookConfiguration
metadata:
  name: networkpolicy-validator
webhooks:
- clientConfig:
    caBundle: "cabundle"
    url: https://127.0.0.1:8443/networkpolicies
  failurePolicy: Fail
  name: networking.k8s-opol.io
  namespaceSelector: {}
  rules:
  - apiGroups:
    - networking.k8s.io
    - extensions
    apiVersions:
    - v1
    - v1beta1
    operations:
    - CREATE
    - UPDATE
    resources:
    - networkpolicies
```
3 - Create your validation rules:
Ex:
```
networkPolicyValidator:
  allowedPolicyTypes: # List with allowed policy types.
  - Egress 
  - Ingress
  podSelector:
    matchLabels:
      rules:
      - name: "LabelCount" 
        operator: "Ge"
        value: 1
  ingress:
    from:
      ipBlock:
        cidr:
          rules:
          - name: "MaskBitsSize"
            operator: "Ge"
            value: 29
        except:
          rules:
          - name: "MaskBitsSize"
            operator: "Ge"
            value: 29
          - name: "ListSize"
            operator: "Le"
            value: 10
      podSelector:
        matchLabels:
          rules:
          - name: "LabelCount"
            operator: "Ge"
            value: "1"
      namespaceSelector:
        matchLabels:
          rules:
          - name: "LabelCount"
            operator: "Ge"
            value: "1"
    ports:
      rules:
      - name: "PortNumber"
        operator: "Ge"
        value: 5000
```
