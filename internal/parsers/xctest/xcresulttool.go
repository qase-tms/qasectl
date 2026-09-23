package xctest

import (
	"errors"
	"fmt"
	"log/slog"
	"os/exec"
	"slices"
	"strings"
)

// legacyExitCode is the status xcresulttool (Xcode 16+) exits with when a legacy
// command is run without --legacy.
const legacyExitCode = 64

// xcresulttool runs `xcrun xcresulttool` with args. Once a call has needed the --legacy
// retry, the flag is passed up front, so later calls no longer pay for a rejected run each.
func (p *Parser) xcresulttool(args ...string) ([]byte, error) {
	if p.legacy {
		return p.runXcresulttool(append(slices.Clone(args), "--legacy"))
	}

	out, err := p.runXcresulttool(args)
	var exitErr *exec.ExitError
	if err == nil || !errors.As(err, &exitErr) || exitErr.ExitCode() != legacyExitCode {
		return out, err
	}

	out, err = p.runXcresulttool(append(slices.Clone(args), "--legacy"))
	if err != nil {
		return nil, fmt.Errorf("with --legacy: %w", err)
	}
	p.legacy = true

	return out, nil
}

// runXcresulttool keeps xcresulttool's stderr in the error: the exit status alone cannot
// tell a missing --legacy flag from a wrong --path.
func (p *Parser) runXcresulttool(args []string) ([]byte, error) {
	const op = "xctest.Parser.runXcresulttool"

	xcrun := p.xcrun
	if xcrun == "" {
		xcrun = "xcrun"
	}

	slog.Debug("executing xcresulttool command", "op", op, "args", args)
	out, err := exec.Command(xcrun, append([]string{"xcresulttool"}, args...)...).Output()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			if stderr := strings.TrimSpace(string(exitErr.Stderr)); stderr != "" {
				return nil, fmt.Errorf("%w: %s", err, stderr)
			}
		}
		return nil, err
	}

	return out, nil
}
