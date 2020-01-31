package quetty

import (
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"
)

type TmuxClient struct {
	SocketPath string
}

func NewTmuxClient() *TmuxClient {
	return &TmuxClient{}
}

func NewTmuxClientWithSocket(socketPath string) *TmuxClient {
	return &TmuxClient{
		SocketPath: socketPath,
	}
}

func (t *TmuxClient) NewSession(sessionName string) error {
	_, err := t.RunCmd("new-session", "-d", "-s", sessionName)
	return err
}

func (t *TmuxClient) NewWindow() error {
	_, err := t.RunCmd("new-window")
	return err
}

func (t *TmuxClient) SplitWindow() error {
	_, err := t.RunCmd("split-window")
	return err
}

func (t *TmuxClient) ListWindows() ([]string, error) {
	listWindowOut, err := t.RunCmd("list-windows", "-F", "#I")
	if err != nil {
		return nil, err
	}
	windows := strings.Split(string(listWindowOut), "\n")
	return windows, nil
}

func (t *TmuxClient) ListPanes() ([]string, error) {
	listPaneOut, err := t.RunCmd("list-panes", "-F", "#D")
	if err != nil {
		return nil, err
	}
	panes := strings.Split(string(listPaneOut), "\n")
	return panes, nil
}

func (t *TmuxClient) ListWindowPanes(win int) ([]string, error) {
	listPaneOut, err := t.RunCmd("list-panes", "-F", "#D", "-t", strconv.Itoa(win))
	if err != nil {
		return nil, err
	}
	panes := strings.Split(string(listPaneOut), "\n")
	return panes, nil
}

func (t *TmuxClient) CurrentWindow() (int, error) {
	curWindowOut, err := t.RunCmd("display-message", "-p", "#I")
	if err != nil {
		return -1, err
	}
	curWindow, err := strconv.Atoi(curWindowOut)
	if err != nil {
		return -1, err
	}
	return curWindow, nil
}

func (t *TmuxClient) CurrentPane() (string, error) {
	curPaneOut, err := t.RunCmd("display-message", "-p", "#D")
	if err != nil {
		return "", err
	}
	return curPaneOut, nil
}

func (t *TmuxClient) TmuxCapturePane(paneId string) (string, error) {
	displayOut, err := t.RunCmd("display-message", "-p", "-t", paneId,
		"#{scroll_region_lower}-#{scroll_position}")
	if err != nil {
		return "", err
	}
	paneVals := strings.Split(displayOut, "-")
	if len(paneVals) != 2 {
		return "", fmt.Errorf("Invalid pane vals")
	}

	capturePaneArgs := []string{"capture-pane", "-p", "-J", "-t", paneId}
	if paneVals[1] != "" {
		scrollHeight, err := strconv.Atoi(paneVals[0])
		if err != nil {
			return "", fmt.Errorf("Couldnt convert scroll height %s",
				paneVals[0])
		}
		scrollPos, err := strconv.Atoi(paneVals[1])
		if err != nil {
			return "", fmt.Errorf("Couldnt convert scroll pos %s",
				paneVals[1])
		}
		bottomPos := scrollHeight - scrollPos
		//scroll position empty implies not in copy mode
		capturePaneArgs = append(capturePaneArgs, "-S", "-"+paneVals[1],
			"-E", strconv.Itoa(bottomPos))
	}
	return t.RunCmd(capturePaneArgs...)
}

func (t *TmuxClient) WriteFromPane(paneId string, writer io.WriteCloser) error {
	defer writer.Close()
	paneContent, err := t.TmuxCapturePane(paneId)
	if err != nil {
		return err
	}

	_, err = io.WriteString(writer, paneContent)
	return err
}

func (t *TmuxClient) RunCmd(args ...string) (string, error) {
	argList := []string{}
	if t.SocketPath != "" {
		argList = append(argList, "-S", t.SocketPath)
	}
	argList = append(argList, args...)
	out, err := RunCmd("tmux", argList...)
	if err != nil {
		return "", fmt.Errorf("Tmux err:%v, out = %s", err, out)
	}
	return out, nil
}

func (t *TmuxClient) LoadBuffer(bufName string, output []byte) error {
	argList := []string{}
	if t.SocketPath != "" {
		argList = append(argList, "-S", t.SocketPath)
	}
	argList = append(argList, "load-buffer", "-b", bufName, "-")
	tmuxCmd := exec.Command("tmux", argList...)
	tmuxInputPipe, err := tmuxCmd.StdinPipe()
	if err != nil {
		return err
	}
	err = tmuxCmd.Start()
	if err != nil {
		return err
	}
	_, err = tmuxInputPipe.Write(output)
	if err != nil {
		return err
	}
	tmuxInputPipe.Close()
	return tmuxCmd.Wait()
}

func (t *TmuxClient) PasteBuffer(bufName string, paneId string) error {
	_, err := t.RunCmd("paste-buffer", "-b", bufName, "-t", paneId)
	return err
}

func (t *TmuxClient) DeleteBuffer(bufName string) error {
	_, err := t.RunCmd("delete-buffer", "-b", bufName)
	return err
}

func (t *TmuxClient) SendString(curPane string, output []byte) error {
	err := t.LoadBuffer("quetty-buffer", output)
	if err != nil {
		return err
	}
	err = t.PasteBuffer("quetty-buffer", curPane)
	if err != nil {
		return err
	}
	err = t.DeleteBuffer("quetty-buffer")
	if err != nil {
		return err
	}
	return nil
}
