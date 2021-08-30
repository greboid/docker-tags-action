package main

import (
	"testing"
)

func Test_getOutput(t *testing.T) {
	type args struct {
		gitRepo         string
		inputRepo       string
		gitRef          string
		gitSHA          string
		inputRegistries string
		fullName        string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "valid input",
			args: args{
				gitRepo:         "group/test",
				inputRepo:       "",
				gitRef:          "refs/heads/master",
				gitSHA:          "abc123",
				inputRegistries: "",
				fullName:        "true",
			},
			want: "::set-output name=tags::docker.io/group/test:dev\n::set-output name=version::abc123",
		},
		{
			name: "valid input no full name",
			args: args{
				gitRepo:         "group/test",
				inputRepo:       "",
				gitRef:          "refs/heads/master",
				gitSHA:          "abc123",
				inputRegistries: "",
				fullName:        "false",
			},
			want: "::set-output name=tags::dev\n::set-output name=version::abc123",
		},
		{
			name: "invalid input",
			args: args{
				gitRepo:         "group/test",
				inputRepo:       "",
				gitRef:          "refs/heads/dev",
				gitSHA:          "abc123",
				inputRegistries: "",
				fullName:        "true",
			},
			want: "::set-output name=tags::\n::set-output name=version::unknown",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getOutput(tt.args.gitRepo, tt.args.inputRepo, tt.args.gitRef, tt.args.gitSHA, tt.args.inputRegistries, defaultSeparator, tt.args.fullName); got != tt.want {
				t.Errorf("getOutput() = %v, want %v", got, tt.want)
			}
		})
	}
}
