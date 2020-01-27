package quetty

import (
	"bytes"
	"os/exec"
)

func RunCmd(name string, arg ...string) ([]byte, error) {
	cmd := exec.Command(name, arg...)
	var output bytes.Buffer
	cmd.Stdout = &output
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	return output.Bytes(), nil
}
