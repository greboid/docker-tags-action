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
	}
	tests := []struct {
		name     string
		args     args
		wantTags []string
	}{
		{
			name:     "single registry",
			args:     args{
				imageName:  "name",
				registries: []string{"registry1"},
				versions:   []string{"1.0.1"},
			},
			wantTags: []string{"registry1/name:1.0.1"},
		},
		{
			name:     "multiple registries",
			args:     args{
				imageName:  "name",
				registries: []string{"registry1", "registry2"},
				versions:   []string{"1.0.1"},
			},
			wantTags: []string{"registry1/name:1.0.1", "registry2/name:1.0.1"},
		},
		{
			name:     "multiple versions",
			args:     args{
				imageName:  "name",
				registries: []string{"registry1"},
				versions:   []string{"1.0.1", "latest"},
			},
			wantTags: []string{"registry1/name:1.0.1", "registry1/name:latest"},
		},
		{
			name:     "multiple both",
			args:     args{
				imageName:  "name",
				registries: []string{"registry1", "registry2"},
				versions:   []string{"1.0.1", "latest"},
			},
			wantTags: []string{"registry1/name:1.0.1", "registry1/name:latest", "registry2/name:1.0.1", "registry2/name:latest"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotTags := getTags(tt.args.imageName, tt.args.registries, tt.args.versions); !reflect.DeepEqual(gotTags, tt.wantTags) {
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