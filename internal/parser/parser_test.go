package parser

import (
	"math"
	"strings"
	"testing"
)

func TestEvalBasicArithmetic(t *testing.T) {
	tests := []struct {
		input string
		want  float64
	}{
		{"2 + 3", 5},
		{"10 - 4", 6},
		{"3 * 7", 21},
		{"20 / 4", 5},
		{"0 + 0", 0},
		{"1.5 + 2.5", 4},
		{"10.5 * 2", 21},
	}
	for _, tt := range tests {
		got, err := Eval(tt.input)
		if err != nil {
			t.Errorf("Eval(%q) error: %v", tt.input, err)
			continue
		}
		if math.Abs(got-tt.want) > 1e-9 {
			t.Errorf("Eval(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestEvalPrecedence(t *testing.T) {
	tests := []struct {
		input string
		want  float64
	}{
		{"2 + 3 * 4", 14},
		{"2 * 3 + 4", 10},
		{"10 - 2 * 3", 4},
		{"10 / 2 + 3", 8},
		{"2 + 6 / 3", 4},
	}
	for _, tt := range tests {
		got, err := Eval(tt.input)
		if err != nil {
			t.Errorf("Eval(%q) error: %v", tt.input, err)
			continue
		}
		if math.Abs(got-tt.want) > 1e-9 {
			t.Errorf("Eval(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestEvalParentheses(t *testing.T) {
	tests := []struct {
		input string
		want  float64
	}{
		{"(2 + 3) * 4", 20},
		{"2 * (3 + 4)", 14},
		{"(10 - 2) * (3 + 1)", 32},
		{"((2 + 3))", 5},
		{"((2 + 3) * (4 - 1))", 15},
	}
	for _, tt := range tests {
		got, err := Eval(tt.input)
		if err != nil {
			t.Errorf("Eval(%q) error: %v", tt.input, err)
			continue
		}
		if math.Abs(got-tt.want) > 1e-9 {
			t.Errorf("Eval(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestEvalUnaryMinus(t *testing.T) {
	tests := []struct {
		input string
		want  float64
	}{
		{"-5", -5},
		{"-5 + 3", -2},
		{"-(2 + 3)", -5},
		{"--5", 5},
	}
	for _, tt := range tests {
		got, err := Eval(tt.input)
		if err != nil {
			t.Errorf("Eval(%q) error: %v", tt.input, err)
			continue
		}
		if math.Abs(got-tt.want) > 1e-9 {
			t.Errorf("Eval(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestEvalWhitespace(t *testing.T) {
	tests := []struct {
		input string
		want  float64
	}{
		{"  2+3  ", 5},
		{"2  +  3", 5},
		{"\t2 + 3\t", 5},
	}
	for _, tt := range tests {
		got, err := Eval(tt.input)
		if err != nil {
			t.Errorf("Eval(%q) error: %v", tt.input, err)
			continue
		}
		if math.Abs(got-tt.want) > 1e-9 {
			t.Errorf("Eval(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestEvalErrors(t *testing.T) {
	tests := []struct {
		input   string
		wantErr string
	}{
		{"", "error: no expression provided"},
		{"   ", "error: no expression provided"},
		{"1 / 0", "error: division by zero"},
		{"(2 + 3", "error: unmatched parenthesis"},
		{"2 + + 3", "error: unexpected token"},
		{"abc", "error: unexpected token"},
		{"2 3", "error: unexpected token"},
	}
	for _, tt := range tests {
		_, err := Eval(tt.input)
		if err == nil {
			t.Errorf("Eval(%q) expected error containing %q, got nil", tt.input, tt.wantErr)
			continue
		}
		if !strings.Contains(err.Error(), tt.wantErr) {
			t.Errorf("Eval(%q) error = %q, want containing %q", tt.input, err.Error(), tt.wantErr)
		}
	}
}
