package main

import "testing"

func Test_parseEnv(t *testing.T) {
	type args struct {
		ruleName string
		env      string
	}
	tests := []struct {
		name       string
		args       args
		wantFirst  string
		wantSecond string
		wantErr    bool
	}{
		{
			name:       "Empty env",
			args:       args{ruleName: "toDev", env: ""},
			wantFirst:  "",
			wantSecond: "",
			wantErr:    true,
		},
		{
			name:       "Rule for found into env",
			args:       args{ruleName: "toDev", env: "dev:/dev/disk2>/hone/user/data.tar.gz;file:/hone/user/data.tar.gz>/dev/disk2"},
			wantFirst:  "",
			wantSecond: "",
			wantErr:    true,
		},
		{
			name:       "Rule 'toFile' env",
			args:       args{ruleName: "toFile", env: "toFile:/dev/disk2>/hone/user/data.tar.gz;toDev:/hone/user/data.tar.gz>/dev/disk2"},
			wantFirst:  "/dev/disk2",
			wantSecond: "/hone/user/data.tar.gz",
		},
		{
			name:       "Rule 'toDev' env",
			args:       args{ruleName: "toDev", env: "toFile:/dev/disk2>/hone/user/data.tar.gz;toDev:/hone/user/data.tar.gz>/dev/disk2"},
			wantFirst:  "/hone/user/data.tar.gz",
			wantSecond: "/dev/disk2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFirst, gotSecond, err := parseEnv(tt.args.ruleName, tt.args.env)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseEnv() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotFirst != tt.wantFirst {
				t.Errorf("parseEnv() gotFirst = %v, want %v", gotFirst, tt.wantFirst)
			}
			if gotSecond != tt.wantSecond {
				t.Errorf("parseEnv() gotSecond = %v, want %v", gotSecond, tt.wantSecond)
			}
		})
	}
}
