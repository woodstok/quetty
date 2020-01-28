package quetty

import (
	"bytes"
	"os/exec"
	"strings"
)

func RunCmd(name string, arg ...string) (string, error) {
	cmd := exec.Command(name, arg...)
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output
	err := cmd.Run()
	outVal := strings.TrimSuffix(string(output.Bytes()), "\n")
	if err != nil {
		return outVal, err
	}
	return outVal, nil
}
