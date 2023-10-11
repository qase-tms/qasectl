package xcresult

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"os/exec"
)

func ToJson(xcResultPath string, fileID *string) (map[string]any, error) {
	args := []string{
		"xcresulttool",
		"get",
		"--path",
		xcResultPath,
		"--format",
		"json",
	}
	if fileID != nil {
		args = append(args, "--id", *fileID)
	}

	cmd := exec.Command(
		"xcrun",
		args...,
	)

	buff, err := cmd.Output()
	if err != nil {
		possErr := string(err.(*exec.ExitError).Stderr)
		fmt.Println("ERROR: ", possErr)

		return nil, errors.Wrap(err, "cmd failed")
	}

	var v map[string]any
	err = json.Unmarshal(buff, &v)
	if err != nil {
		return nil, errors.Wrap(err, "json unmarshal")
	}

	return v, nil
}

func OpenAttachment(xcResultPath string, fileID *string) ([]byte, error) {
	args := []string{
		"xcresulttool",
		"get",
		"--path",
		xcResultPath,
		"--format",
		"raw",
	}
	if fileID != nil {
		args = append(args, "--id", *fileID)
	}

	cmd := exec.Command(
		"xcrun",
		args...,
	)

	return cmd.Output()
}
