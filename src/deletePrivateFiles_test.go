package main

import (
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"
)

func Test_PrivateRegex(t *testing.T) {
	tests := []struct {
		name        string
		pathtoCheck string
		match       bool
	}{
		{"Code file", "foo.cpp", false},
		{"header", "foo.h", false},
		{"uasset", "foo.uasset", false},
		{"private cpp file", "foo.cpp.private.0", true},
		{"private h file", "foo.h.private.1", true},
		{"private uasset file", "foo.uasset.private.2", true},
		{"private file 11 (not a match)", "foo.uasset.private.11", false},
		{"private not suffix", "foo.private.cpp", false},
		{"private folder", "private/foo.cpp", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if match := rx.MatchString(tt.pathtoCheck); match != tt.match {
				t.Errorf("testRegexMatch() = %v, want %v", match, tt.match)
			}
		})
	}
}

func Test_findPrivateFiles(t *testing.T) {
	tests := []struct {
		name          string
		filesToCreate []string
		want          []string
	}{
		{"No Private Files", []string{"a.cpp", "b.h"}, []string{}},
		{"Private Files", []string{"a.cpp.private.0", "b.h.private.1"}, []string{"a.cpp.private.0", "b.h.private.1"}},
		{"Mixed Files", []string{"a.h.private.9", "b.cpp"}, []string{"a.h.private.9"}},
		{"binary Files", []string{"foo.uasset", "bar.fbx"}, []string{}},
		{"No numeric suffix Files", []string{"foo.private"}, []string{}},
		{"Private Binary Files", []string{"foo.uasset.private.3", "bar.fbx.private.4"}, []string{"foo.uasset.private.3", "bar.fbx.private.4"}},
		{"Private not suffix", []string{"foo.private.cpp"}, []string{}},
		{"subdir", []string{"foo/a.cpp", "foo/b.cpp.private.5"}, []string{"foo/b.cpp.private.5"}},
		{"Private subdir", []string{"private/foo.cpp"}, []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDir := t.TempDir()
			for _, toCreate := range tt.filesToCreate {
				absPath := filepath.Join(testDir, toCreate)
				basePath := filepath.Dir(absPath)
				if err := os.MkdirAll(basePath, 0777); err != nil {
					t.Fatalf("Couldnt create subdir: %v", err)
				}
				f, err := os.Create(absPath)
				if err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
				_ = f.Close()
			}

			got := findPrivateFiles(testDir)
			sort.Strings(got)
			sort.Strings(tt.want)
			if got := findPrivateFiles(testDir); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findPrivateFiles() = %v, want %v", got, tt.want)
			}
		})
	}
}
