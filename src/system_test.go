package quetty

import "testing"

func TestRunCmd(t *testing.T) {
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
