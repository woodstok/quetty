package quetty

import (
	"os"
	"testing"
)

const TmpSocketPath = "/tmp/tmpTmuxSocket"

func TestListPanes(t *testing.T) {

	defer os.Remove(TmpSocketPath)

	// Start tmux session

	tmuxClient := NewTmuxClientWithSocket(TmpSocketPath)
	err := tmuxClient.NewSession("temp-session")
	assertNoError(t, err)

	err := tmuxClient.NewWindow()
	assertNoError(t, err)

	err := tmuxClient.SplitWindow()
	assertNoError(t, err)
	err := tmuxClient.SplitWindow()
	assertNoError(t, err)

	t.Run("Basic command", func(t *testing.T) {
		got, err := RunCmd("echo", "hi")
		assertNoError(t, err)
		assertStringEqual(t, "hi\n", string(got))
	})
	t.Run("Nonexistent command", func(t *testing.T) {
		_, err := RunCmd("thiscommanddoesnotexist", "hi")
		assertSomeError(t, err)
	})
}
