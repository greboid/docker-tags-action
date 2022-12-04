package main

import (
	"reflect"
	"testing"

	"github.com/blang/semver/v4"
)

func Test_getOutput(t *testing.T) {
	type args struct {
		gitRepo         string
		inputRepo       string
		gitRef          string
		gitSHA          string
		inputRegistries string
		fullName        string
		latestVersion   *semver.Version
	}
	tests := []struct {
		name string
		args args
		want map[string]string
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
			want: map[string]string{"tags": "docker.io/group/test:dev", "version": "abc123"},
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
			want: map[string]string{"tags": "dev", "version": "abc123"},
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
			want: map[string]string{"tags": "", "version": "unknown"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getOutput(tt.args.gitRepo, tt.args.inputRepo, tt.args.gitRef, tt.args.gitSHA, tt.args.inputRegistries, defaultSeparator, tt.args.fullName, tt.args.latestVersion); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getOutput() = %v, want %v", got, tt.want)
			}
		})
	}
}
