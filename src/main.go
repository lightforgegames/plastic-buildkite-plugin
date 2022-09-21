package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

func Shellout(command string) (error, string, string) {
	const ShellToUse = "bash"
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command(ShellToUse, "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return err, stdout.String(), stderr.String()
}

func get_comment(changeset int, branch string) (string, error) {
	query := ""
	if changeset == -1 {
		query = fmt.Sprintf("find changeset \"where branch = '%s'\" --format=\"{comment}\" order by changesetId desc LIMIT 1 --nototal", branch)
	} else {
		query = fmt.Sprintf("log cs:%d --csformat={comment}", changeset)
	}
	err, out, _ := Shellout("cm " + query)
	return out, err
}

func main() {
	cd, _ := os.Getwd()
	fmt.Println("go: executing plastic plugin from " + cd)
	_, err := exec.Command("cm", "ss").CombinedOutput()
	if err != nil {
		// cm ss failed, so we can set up a workspace here now.
		repoPath := os.Getenv("BUILDKITE_REPO")
		pipelineName := os.Getenv("BUILDKITE_PIPELINE_NAME")
		agentId := os.Getenv("BUILDKITE_AGENT_ID")
		workspaceName := fmt.Sprintf("buildkite-%s-%s", pipelineName, agentId)
		out, err := exec.Command("cm", "workspace", "create", workspaceName, ".", "--repository="+repoPath).CombinedOutput()
		if err != nil {
			fmt.Printf("Failed to create workspace `%s`: %v.\n%s", workspaceName, err, string(out))
			os.Exit(1)
		}
	} else {
		fmt.Println("Cleaning workspace...")
		if out, err := exec.Command("cm", "undo", ".", "-R").CombinedOutput(); err != nil {
			fmt.Printf("Failed to undo changes: : %v.\n%s", err, string(out))
			os.Exit(1)
		}
	}

	// figure out our target branch and changeset.
	// Start with the branch. if the branch has been specified, then use that. Let it be overridden by a cs though.
	branch := os.Getenv("BUILDKITE_BRANCH")
	target := "br:" + branch

	changeset := -1

	revision := os.Getenv("BUILDKITE_COMMIT")
	if !(revision == "" || revision == "HEAD") {
		// If the revision isn't empty, or head, then set the target to the specified changeset
		if changeset, err = strconv.Atoi(revision); err != nil || changeset < 1 {
			fmt.Printf("Invalid changeset specified. Expected a numeric value but got `%s`", revision)
			os.Exit(1)
		}
	}

	if len(target) == 3 {
		fmt.Printf("Invalid target, expected either a branch or a changeset but got `%s`", target)
		os.Exit(1)
	}

	fmt.Println("Setting workspace to " + target)
	if out, err := exec.Command("cm", "switch", target).CombinedOutput(); err != nil {
		fmt.Printf("Failed to update workspace: : %v.\n%s", err, string(out))
		os.Exit(1)
	}

	// Finally, set metadata
	comment, err := get_comment(changeset, branch)
	if err != nil {
		fmt.Printf("Failed to get comment for `%v:%s`", changeset, branch)
		fmt.Printf("Failed toget comment: : %v.\n%s", err, comment)
		os.Exit(1)
	}
	commitMetadata := fmt.Sprintf("commit %s\n\n\t%s", revision, comment)
	if out, err := exec.Command("buildkite-agent", "meta-data", "set", "buildkite:git:commit", commitMetadata).CombinedOutput(); err != nil {
		fmt.Printf("Failed to set metadata: : %v.\n%s", err, string(out))
		os.Exit(1)
	}
}
