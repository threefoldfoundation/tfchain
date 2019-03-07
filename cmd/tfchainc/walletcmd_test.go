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
			args: `\x00\x01\x02`,
			want: string([]byte{0, 1, 2}),
		}, {
			name: "2",
			args: `\x00\x05\x13`,
			want: string([]byte{0, 5, 19}),
		}, {
			name: "3",
			args: `\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00`,
			want: string([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}),
		}, {
			name: "4",
			args: "Hello, world!",
			want: "Hello, world!",
		}, {
			name: "5",
			args: "რეგისტრაცია Unicode-ის მეათე",
			want: "რეგისტრაცია Unicode-ის მეათე",
		}, {
			name: "6",
			args: `\u10e0\u10d4\u10d2\u10d8\u10e1\u10e2\u10e0\u10d0\u10ea\u10d8\u10d0 Unicode-\u10d8\u10e1 \u10db\u10d4\u10d0\u10d7\u10d4`,
			want: "რეგისტრაცია Unicode-ის მეათე",
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
