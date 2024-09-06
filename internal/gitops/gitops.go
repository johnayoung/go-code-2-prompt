package gitops

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// GetStagedDiff returns the diff of staged changes
func GetStagedDiff(repoPath string) (string, error) {
	cmd := exec.Command("git", "-C", repoPath, "diff", "--staged")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("error getting staged diff: %v", err)
	}
	return out.String(), nil
}

// GetBranchDiff returns the diff between two branches
func GetBranchDiff(repoPath, branch1, branch2 string) (string, error) {
	cmd := exec.Command("git", "-C", repoPath, "diff", branch1+".."+branch2)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("error getting branch diff: %v", err)
	}
	return out.String(), nil
}

// GetGitLog returns the git log between two branches
func GetGitLog(repoPath, branch1, branch2 string) (string, error) {
	cmd := exec.Command("git", "-C", repoPath, "log", "--oneline", branch1+".."+branch2)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("error getting git log: %v", err)
	}
	return out.String(), nil
}

// IsGitRepository checks if the given path is a git repository
func IsGitRepository(path string) bool {
	cmd := exec.Command("git", "-C", path, "rev-parse", "--is-inside-work-tree")
	err := cmd.Run()
	return err == nil
}

// GetCurrentBranch returns the name of the current git branch
func GetCurrentBranch(repoPath string) (string, error) {
	cmd := exec.Command("git", "-C", repoPath, "rev-parse", "--abbrev-ref", "HEAD")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("error getting current branch: %v", err)
	}
	return strings.TrimSpace(out.String()), nil
}
