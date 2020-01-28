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

	err = tmuxClient.NewWindow()
	assertNoError(t, err)
	err = tmuxClient.NewWindow()
	assertNoError(t, err)

	err = tmuxClient.SplitWindow()
	assertNoError(t, err)
	err = tmuxClient.SplitWindow()
	assertNoError(t, err)

	t.Run("List Windows", func(t *testing.T) {
		windows, err := tmuxClient.ListWindows()
		expectedWindows := 2
		assertNoError(t, err)
		if len(windows) != expectedWindows {
			t.Errorf("Got %d windows, expected %d windows",
				len(windows), expectedWindows)
		}
	})

	t.Run("List Panes", func(t *testing.T) {
		panes, err := tmuxClient.ListPanes()
		expectedPanes := 3
		assertNoError(t, err)
		if len(panes) != expectedPanes {
			t.Errorf("Got %d panes, expected %d panes",
				len(panes), expectedPanes)
		}
	})

	t.Run("List Panes of a window", func(t *testing.T) {
		panes, err := tmuxClient.ListWindowPanes(2)
		expectedPanes := 1
		assertNoError(t, err)
		if len(panes) != expectedPanes {
			t.Errorf("Got %d panes, expected %d panes",
				len(panes), expectedPanes)
		}
	})
}

func TestCapturePane(t *testing.T) {

	// Temporarily test current tmux session
	// change to
	// create session
	// create window with output
	// capture pane and test
	// Start tmux session

	tmuxClient := NewTmuxClient()

	_, err := tmuxClient.TmuxCapturePane("%23")
	assertNoError(t, err)

	// t.Run("List Windows", func(t *testing.T) {
	// 	windows, err := tmuxClient.ListWindows()
	// 	expectedWindows := 2
	// 	assertNoError(t, err)
	// 	if len(windows) != expectedWindows {
	// 		t.Errorf("Got %d windows, expected %d windows",
	// 			len(windows), expectedWindows)
	// 	}
	// })

	// t.Run("List Panes", func(t *testing.T) {
	// 	panes, err := tmuxClient.ListPanes()
	// 	expectedPanes := 3
	// 	assertNoError(t, err)
	// 	if len(panes) != expectedPanes {
	// 		t.Errorf("Got %d panes, expected %d panes",
	// 			len(panes), expectedPanes)
	// 	}
	// })

	// t.Run("List Panes of a window", func(t *testing.T) {
	// 	panes, err := tmuxClient.ListWindowPanes(2)
	// 	expectedPanes := 1
	// 	assertNoError(t, err)
	// 	if len(panes) != expectedPanes {
	// 		t.Errorf("Got %d panes, expected %d panes",
	// 			len(panes), expectedPanes)
	// 	}
	// })
}
