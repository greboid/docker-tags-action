package main

import (
	"reflect"
	"testing"
)

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

func Test_refToVersions(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		wantVersions []string
	}{
		{
			name:         "v1.0.0",
			input:        "refs/tags/v1.0.0",
			wantVersions: []string{"1.0.0", "1.0", "1"},
		},
		{
			name:         "1.0.0",
			input:        "refs/tags/v1.0.0",
			wantVersions: []string{"1.0.0", "1.0", "1"},
		},
		{
			name:         "v1.0.0.0.0.0",
			input:        "refs/tags/v1.0.0.0.0.0",
			wantVersions: nil,
		},
		{
			name:         "1.0",
			input:        "refs/tags/1.0",
			wantVersions: nil,
		},
		{
			name:         "v1.0",
			input:        "refs/tags/v1.0",
			wantVersions: nil,
		},
		{
			name:         "1",
			input:        "refs/tags/1",
			wantVersions: nil,
		},
		{
			name:         "v1",
			input:        "refs/tags/v1",
			wantVersions: nil,
		},
		{
			name:         "master",
			input:        "refs/heads/master",
			wantVersions: []string{"latest"},
		},
		{
			name:         "main",
			input:        "refs/heads/main",
			wantVersions: []string{"latest"},
		},
		{
			name:         "non numeric number",
			input:        "refs/tags/v1.a",
			wantVersions: nil,
		},
		{
			name:         "invalid ref",
			input:        "refs/heads/dev",
			wantVersions: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotVersions := refToVersions(tt.input); !reflect.DeepEqual(gotVersions, tt.wantVersions) {
				t.Errorf("refToVersions() = %v, want %v", gotVersions, tt.wantVersions)
			}
		})
	}
}

func Test_getTags(t *testing.T) {
	type args struct {
		imageName  string
		registries []string
		versions   []string
		fullname bool
	}
	tests := []struct {
		name     string
		args     args
		wantTags []string
	}{
		{
			name: "single registry (full name)",
			args: args{
				imageName:  "name",
				registries: []string{"registry1"},
				versions:   []string{"1.0.1"},
				fullname: true,
			},
			wantTags: []string{"registry1/name:1.0.1"},
		},
		{
			name: "multiple registries (full name)",
			args: args{
				imageName:  "name",
				registries: []string{"registry1", "registry2"},
				versions:   []string{"1.0.1"},
				fullname: true,
			},
			wantTags: []string{"registry1/name:1.0.1", "registry2/name:1.0.1"},
		},
		{
			name: "multiple versions (full name)",
			args: args{
				imageName:  "name",
				registries: []string{"registry1"},
				versions:   []string{"1.0.1", "latest"},
				fullname: true,
			},
			wantTags: []string{"registry1/name:1.0.1", "registry1/name:latest"},
		},
		{
			name: "multiple versions (mo full name)",
			args: args{
				imageName:  "name",
				registries: []string{"registry1"},
				versions:   []string{"1.0.1", "latest"},
				fullname: false,
			},
			wantTags: []string{"1.0.1", "latest"},
		},
		{
			name: "multiple both (full name)",
			args: args{
				imageName:  "name",
				registries: []string{"registry1", "registry2"},
				versions:   []string{"1.0.1", "latest"},
				fullname: true,
			},
			wantTags: []string{"registry1/name:1.0.1", "registry1/name:latest", "registry2/name:1.0.1", "registry2/name:latest"},
		},
		{
			name: "multiple both (no full name)",
			args: args{
				imageName:  "name",
				registries: []string{"registry1", "registry2"},
				versions:   []string{"1.0.1", "latest"},
				fullname: false,
			},
			wantTags: []string{"1.0.1", "latest", "1.0.1", "latest"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotTags := getTags(tt.args.imageName, tt.args.registries, tt.args.versions, tt.args.fullname); !reflect.DeepEqual(gotTags, tt.wantTags) {
				t.Errorf("getTags() = %v, want %v", gotTags, tt.wantTags)
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

func Test_getOutput(t *testing.T) {
	type args struct {
		gitRepo         string
		inputRepo       string
		gitRef          string
		gitSHA          string
		inputRegistries string
		fullName	string
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
				fullName: "true",
			},
			want: "::set-output name=tags::docker.io/group/test:latest\n::set-output name=version::abc123",
		},
		{
			name: "valid input no full name",
			args: args{
				gitRepo:         "group/test",
				inputRepo:       "",
				gitRef:          "refs/heads/master",
				gitSHA:          "abc123",
				inputRegistries: "",
				fullName: "false",
			},
			want: "::set-output name=tags::latest\n::set-output name=version::abc123",
		},
		{
			name: "invalid input",
			args: args{
				gitRepo:         "group/test",
				inputRepo:       "",
				gitRef:          "refs/heads/dev",
				gitSHA:          "abc123",
				inputRegistries: "",
				fullName: "true",
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

func Test_refToVersion(t *testing.T) {
	type args struct {
		ref string
		sha string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "v1.0.0",
			args: args{
				ref: "refs/tags/v1.0.0",
				sha: "abc123",
			},
			want: "1.0.0",
		},
		{
			name: "1.0.0",
			args: args{
				ref: "refs/tags/1.0.0",
				sha: "abc123",
			},
			want: "1.0.0",
		},
		{
			name: "v1.0.0.0.0.0",
			args: args{
				ref: "refs/tags/v1.0.0.0.0.0",
				sha: "abc123",
			},
			want: "unknown",
		},
		{
			name: "1.0",
			args: args{
				ref: "refs/tags/1.0",
				sha: "abc123",
			},
			want: "unknown",
		},
		{
			name: "v1.0",
			args: args{
				ref: "refs/tags/v1.0",
				sha: "abc123",
			},
			want: "unknown",
		},
		{
			name: "v1",
			args: args{
				ref: "refs/tags/v1",
				sha: "abc123",
			},
			want: "unknown",
		},
		{
			name: "1",
			args: args{
				ref: "refs/tags/1",
				sha: "abc123",
			},
			want: "unknown",
		},
		{
			name: "master",
			args: args{
				ref: "refs/heads/master",
				sha: "abc123",
			},
			want: "abc123",
		},
		{
			name: "main",
			args: args{
				ref: "refs/heads/main",
				sha: "abc123",
			},
			want: "abc123",
		},
		{
			name: "non numeric number",
			args: args{
				ref: "refs/tags/v1.a",
				sha: "abc123",
			},
			want: "unknown",
		},
		{
			name: "invalid ref",
			args: args{
				ref: "refs/heads/dev",
				sha: "abc123",
			},
			want: "unknown",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := refToVersion(tt.args.ref, tt.args.sha); got != tt.want {
				t.Errorf("refToVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
