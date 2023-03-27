package main

import "testing"

func Test_friendlyBranchName(t *testing.T) {
	tests := []struct {
		name       string
		branchName string
		want       string
		wantErr    bool
	}{
		{"/main to main", "/main", "main", false},
		{"no leading slash", "main", "", true},
		{"no trailing slash", "main/", "", true},
		{"/main/child to main-child", "/main/child", "main-child", false},
		{"main/child errors", "main/child", "", true},
		{"/main/child_branch", "/main/child_branch", "main-child_branch", false},
		{"main/child_branch", "main/child_branch", "", true},

		{"invalid characters in branch", "main/with-dashes", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getFriendlyBranchName(tt.branchName)
			if (err != nil) != tt.wantErr {
				t.Errorf("getFriendlyBranchName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getFriendlyBranchName() got = %v, want %v", got, tt.want)
			}
		})
	}
}
