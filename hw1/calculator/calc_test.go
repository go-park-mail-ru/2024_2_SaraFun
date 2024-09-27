package main

import "testing"

func TestEvaluateExpression(t *testing.T) {
	tests := []struct {
		expr      string
		expected  float64
		expectErr bool
	}{
		{"3 + 5", 8, false},
		{"10 - 2", 8, false},
		{"4 * 2", 8, false},
		{"16 / 2", 8, false},
		{"3 + 4 * 2", 11, false},
		{"(1 + 2) * 3", 9, false},
		{"10 / (2 + 3)", 2, false},
		{"(10 - 5) * (2 + 3)", 25, false},
		{"5 + (3 * (2 - 1))", 8, false},
		{"3.5 + 2.5", 6, false},
		{"7 + 1 * (-1)", 6, false},
		{"7 + (-1) * -1", 8, false},
		{"7 + -1 * -1", 8, false},
		{"(90 - 91 + 1) * 3", 0, false},
		{"(90 - 91 + 1) / 3", 0, false},
		{"(90 - 91) * 3", -3, false},
		{"(90 - 91) / 3", -0.3333, false},
		{"0", 0, false},
		{"48", 48, false},
		{"-52", -52, false},
		{"-8", -8, false},
		{"0-8", -8, false},
		{"5*-8", -40, false},
		{"-(90 - 91 + 1) / 3", 0, false},
		{"-1 + (90 - 91 + 1) / 3", -1, false},
		{"-(-11-(1*20/2)-11/2*3)", 37.5, false},
		{"() + 83", 0, true},           // Некорректное выражение
		{"52.3.2 + 8", 0, true},        // Некорректное представление числа
		{"10 / 0", 0, true},            // Деление на ноль
		{"3 / (90 - 89 - 1)", 0, true}, // Деление на ноль
		{"3 +", 0, true},               // Некорректное выражение
		{"* 3 + 5", 0, true},           // Некорректное выражение
		{"/ 3 + 5", 0, true},           // Некорректное выражение
		{"5 + 245 - ()*8", 0, true},    // Некорректное выражение
		{"--911", 0, true},             // Некорректное выражение
		{"*", 0, true},                 // Некорректное выражение
		{"+-", 0, true},                // Некорректное выражение
		{"-", 0, true},                 // Некорректное выражение
	}

	for _, test := range tests {
		t.Run(test.expr, func(t *testing.T) {
			result, err := calcExpr(test.expr)
			if test.expectErr {
				if err == nil {
					t.Errorf("expected an error for expression %q but got none", test.expr)
				}
			} else {
				if err != nil {
					t.Errorf("did not expect an error for expression %q but got: %v", test.expr, err)
				}
				if result != test.expected {
					t.Errorf("for expression %q expected %v but got %v", test.expr, test.expected, result)
				}
			}
		})
	}
}
