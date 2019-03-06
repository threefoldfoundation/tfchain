package main

import "testing"

func Test_parseDescription(t *testing.T) {
	tests := []struct {
		name string
		args string
		want string
	}{
		{
			name: "1",
			args: "\\x00\\x01\\x02",
			want: string([]byte{0, 1, 2}),
		}, {
			name: "2",
			args: "\\x00\\x05\\x13",
			want: string([]byte{0, 5, 19}),
		}, {
			name: "3",
			args: "\\x00\\x00\\x00\\x00\\x00\\x00\\x00\\x00\\x00\\x00",
			want: string([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseDescription(tt.args); got != tt.want {
				t.Errorf("parseDescription() = %v, want %v", got, tt.want)
			}
		})
	}
}
