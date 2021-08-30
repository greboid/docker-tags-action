package main

import (
	"reflect"
	"testing"
)

func Test_getSeparator(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "Blank input",
			input: "",
			want:  defaultSeparator,
		},
		{
			name:  "Comma",
			input: ",",
			want:  ",",
		},
		{
			name:  "Space",
			input: " ",
			want:  " ",
		},
		{
			name:  "Multiple",
			input: "test",
			want:  "test",
		},
		{
			name:  "Unicode",
			input: "👩🏿‍🚀",
			want:  "👩🏿‍🚀",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getSeparator(tt.input); got != tt.want {
				t.Errorf("getSeparator() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getFullName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "Blank",
			input: "",
			want:  true,
		},
		{
			name:  "true",
			input: "true",
			want:  true,
		},
		{
			name:  "1",
			input: "1",
			want:  true,
		},
		{
			name:  "0",
			input: "0",
			want:  false,
		},
		{
			name:  "false",
			input: "false",
			want:  false,
		},
		{
			name:  "random",
			input: "random",
			want:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getFullName(tt.input); got != tt.want {
				t.Errorf("getFullName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getImageName(t *testing.T) {
	type args struct {
		gitRepo   string
		inputRepo string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "git repo only",
			args: args{
				gitRepo:   "git",
				inputRepo: "",
			},
			want: "git",
		},
		{
			name: "input repo provided",
			args: args{
				gitRepo:   "git",
				inputRepo: "image",
			},
			want: "image",
		},
		{
			name: "trimmed input repo",
			args: args{
				gitRepo:   "git",
				inputRepo: "",
			},
			want: "git",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getImageName(tt.args.gitRepo, tt.args.inputRepo); got != tt.want {
				t.Errorf("getImageName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseRegistriesInput(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "Empty",
			input: "",
			want:  []string{"docker.io"},
		},
		{
			name:  "Empty with whitespace",
			input: "    ",
			want:  []string{"docker.io"},
		},
		{
			name:  "One",
			input: "reg1",
			want:  []string{"reg1"},
		},
		{
			name:  "One with whitespace",
			input: " reg1 ",
			want:  []string{"reg1"},
		},
		{
			name:  "Multiple",
			input: "reg1, reg2",
			want:  []string{"reg1", "reg2"},
		},
		{
			name:  "Multiple with empty",
			input: "reg1, ,reg2",
			want:  []string{"reg1", "reg2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseRegistriesInput(tt.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseRegistriesInput() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_splitRepo(t *testing.T) {
	tests := []struct {
		name    string
		repo    string
		want    string
		want1   string
		wantErr bool
	}{
		{
			name:    "Normal",
			repo:    "owner/repo",
			want:    "owner",
			want1:   "repo",
			wantErr: false,
		},
		{
			name:    "Not enough parts",
			repo:    "owner",
			want:    "",
			want1:   "",
			wantErr: true,
		},
		{
			name:    "Too many parts",
			repo:    "owner/repo/test",
			want:    "",
			want1:   "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := splitRepo(tt.repo)
			if (err != nil) != tt.wantErr {
				t.Errorf("splitRepo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("splitRepo() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("splitRepo() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
