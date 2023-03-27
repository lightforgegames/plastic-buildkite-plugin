package main

import "os/exec"

func setMetadata(name string, value string) (string, error) {
	out, err := exec.Command("buildkite-agent", "meta-data", "set", name, value).CombinedOutput()
	return string(out), err
}

func getMetadata(name string, defaultValue string) (string, error) {
	out, err := exec.Command("buildkite-agent", "meta-data", "get", name, "--default", defaultValue).CombinedOutput()
	if err == nil {
		return string(out), err
	} else {
		return "", err
	}
}

func annotate(style, context, message string) error {
	_, err := exec.Command("buildkite-agent", "annotate", message, "--context", context, "--style", style).CombinedOutput()
	return err
}
