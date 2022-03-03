package main

import "testing"

func Test_myTestableFunction(t *testing.T) {
	tests := []struct {
		name  string
		value int
		want  int
	}{
		{
			name:  "Simple test 1",
			value: 1000,
			want:  1001,
		},
		{
			name:  "Simple test 2",
			value: 0,
			want:  1,
		},
		{
			name:  "Simple test 3",
			value: 199,
			want:  200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := myTestableFunction(tt.value); got != tt.want {
				t.Errorf("myTestableFunction() = %v, want %v", got, tt.want)
			}
		})
	}
}
