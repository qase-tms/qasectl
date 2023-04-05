package xcresult

import (
	"os/exec"
)

func ToJson(xcResultPath string) ([]byte, error) {
	/**
	todo support:
	if file_id is not None:
	           prams.extend(["--id", file_id])
	*/
	// xcrun xcresulttool get --path folder.xcresult --format json
	cmd := exec.Command(
		"xcrun",
		"xcresulttool",
		"get",
		"--path",
		xcResultPath,
		"--format",
		"json",
	)

	return cmd.Output()
}
