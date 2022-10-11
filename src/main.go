package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func Shellout(command string) (string, string, error) {
	const ShellToUse = "bash"
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command(ShellToUse, "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

func getHead(branch string) (string, error) {
	out, _, err := Shellout(fmt.Sprintf("find changeset \"where branch = '%s'\" --format=\"{changesetid}\" order by changesetId desc LIMIT 1 --nototal", branch))
	return out, err
}

func getComment(changeset int) (string, error) {
	out, _, err := Shellout(fmt.Sprintf("cm log cs:%d --csformat={comment}", changeset))
	return out, err
}

func main() {
	appDataPath := os.Getenv("LocalAppData")
	fmt.Printf("Using Appdata: %q\n", appDataPath)

	cd, _ := os.Getwd()

	fmt.Println("go: executing plastic plugin from " + cd)
	_, err := os.Stat(".plastic/plastic.selector")
	selectorString := ""
	if bHasSelector := err == nil; bHasSelector {
		selectorString = "--selector=.plastic/plastic.selector"
	}

	repoPath := os.Getenv("BUILDKITE_REPO")
	pipelineName := os.Getenv("BUILDKITE_PIPELINE_NAME")

	workspaceName, found := os.LookupEnv("BUILDKITE_PLUGIN_PLASTIC_WORKSPACENAME")
	if !found {
		workspaceName = fmt.Sprintf("buildkite-%s", pipelineName)
	}

	fmt.Printf("Creating workspace %q for repository %q\n", workspaceName, repoPath)
	out, err := exec.Command("cm", "workspace", "create", workspaceName, ".", selectorString).CombinedOutput()
	if err != nil {
		if !strings.Contains(string(out), "already exists.") {
			fmt.Printf("Failed to create workspace `%s`: %v.\n%s", workspaceName, err, string(out))
			os.Exit(1)
		}
	}

	fmt.Println("Cleaning workspace of any changes...")
	if out, err := exec.Command("cm", "undo", ".", "-R").CombinedOutput(); err != nil {
		fmt.Printf("Failed to undo changes: : %v.\n%s", err, string(out))
		os.Exit(1)
	}

	// figure out our target branch and changeset.
	// Start with the branch. if the branch has been specified, then use that. Let it be overridden by a cs though.
	branch := os.Getenv("BUILDKITE_BRANCH")
	target := "br:" + branch

	changeset := -1

	revision := os.Getenv("BUILDKITE_COMMIT")
	if revision == "" || revision == "HEAD" {
		revision, err = getHead(branch)
		if err != nil {
			fmt.Printf("failed to get head of branch %q: %v", branch, err)
			os.Exit(1)
		}
	} else {
		target = "cs:" + revision
	}

	if changeset, err = strconv.Atoi(revision); err != nil || changeset < 1 {
		fmt.Printf("Invalid changeset specified. Expected a numeric value but got `%s`\n", revision)
		os.Exit(1)
	}

	// If the revision isn't empty, or head, then set the target to the specified changeset
	if len(target) == 3 {
		fmt.Printf("Invalid target, expected either a branch or a changeset but got `%s`\n", target)
		os.Exit(1)
	}

	// Set metadata before updating, as updating can take minutes.
	comment, err := getComment(changeset)
	if err != nil {
		fmt.Printf("Failed to get comment for `%v:%s`\n", changeset, branch)
		fmt.Printf("Failed to get comment: : %v.\n%s\n", err, comment)
		os.Exit(1)
	}
	commitMetadata := fmt.Sprintf("commit %s\n\n\t%s", revision, comment)
	if out, err := exec.Command("buildkite-agent", "meta-data", "set", "buildkite:git:commit", commitMetadata).CombinedOutput(); err != nil {
		fmt.Printf("Failed to set metadata: : %v.\n%s\n", err, string(out))
		os.Exit(1)
	}

	fmt.Println("Setting workspace to " + target)
	if out, err := exec.Command("cm", "switch", target).CombinedOutput(); err != nil {
		fmt.Printf("Failed to update workspace: : %v.\n%s\n", err, string(out))
		os.Exit(1)
	}

	fmt.Println("Update complete.")
}
