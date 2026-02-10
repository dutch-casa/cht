package main

import (
	"errors"
	"os"
	"os/exec"
	"strings"
)

func fzfSelect() (string, error) {
	topicFile, err := ensureCache()
	if err != nil {
		return "", err
	}

	f, err := os.Open(topicFile)
	if err != nil {
		return "", err
	}
	defer func() { _ = f.Close() }()

	fzf := exec.Command("fzf",
		"--preview", "curl -s cheat.sh/{}",
		"--preview-window", "right:70%:wrap",
	)
	fzf.Stdin = f
	fzf.Stderr = os.Stderr

	out, err := fzf.Output()
	if err != nil {
		// fzf exits 130 on Ctrl-C / Esc — not an error.
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) && exitErr.ExitCode() == 130 {
			return "", nil
		}
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}
