package somesql

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFieldValueFunctions(t *testing.T) {
	type args struct {
		funcs []string
	}
	tests := []struct {
		name  string
		args  args
		field string
		value string
	}{
		{
			name:  "Zero",
			args:  args{},
			field: "",
			value: "",
		},
		{
			name:  "One",
			args:  args{funcs: []string{"FOO"}},
			field: "FOO",
			value: "",
		},
		{
			name:  "Two",
			args:  args{funcs: []string{"FOO", "BAR"}},
			field: "FOO",
			value: "BAR",
		},
		{
			name:  "More",
			args:  args{funcs: []string{"FOO", "BAR", "BAZ", "BOINK"}},
			field: "FOO",
			value: "BAR",
		},
	}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field, value := getFieldValueFunctions(tt.args.funcs)
			assert.Equal(t, tt.field, field, fmt.Sprintf("%d: Field function invalid", i))
			assert.Equal(t, tt.value, value, fmt.Sprintf("%d: Value function invalid", i))
		})
	}
}
