package admission

import "testing"

func TestNumberOperators(t *testing.T) {

	tests := []struct {
		x        int
		op       Operator
		y        int
		expected bool
	}{
		{
			1,
			"Equals",
			2,
			false,
		},
		{
			1,
			"Lt",
			2,
			true,
		},
		{
			2,
			"Ge",
			1,
			true,
		},
		{
			2,
			"Gt",
			1,
			true,
		},
		{
			2,
			"Le",
			1,
			false,
		},
	}

	for _, i := range tests {

		result, _ := operatorExec(i.x, i.y, i.op)

		if result != i.expected {
			t.Errorf(" %v", result)
		}

	}

}
