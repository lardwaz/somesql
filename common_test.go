package somesql

import "testing"

func Test_getFieldValueFunctions(t *testing.T) {
	type args struct {
		funcs []string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
	}{
		{
			name:  "Zero",
			args:  args{},
			want:  "",
			want1: "",
		},
		{
			name:  "One",
			args:  args{funcs: []string{"FOO"}},
			want:  "FOO",
			want1: "",
		},
		{
			name:  "Two",
			args:  args{funcs: []string{"FOO", "BAR"}},
			want:  "FOO",
			want1: "BAR",
		},
		{
			name:  "More",
			args:  args{funcs: []string{"FOO", "BAR", "BAZ", "BOINK"}},
			want:  "FOO",
			want1: "BAR",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := getFieldValueFunctions(tt.args.funcs)
			if got != tt.want {
				t.Errorf("getFieldValueFunctions() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("getFieldValueFunctions() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
