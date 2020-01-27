package quetty

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
	return t.RunCmd("new-session", sessionName)
}

func (t *TmuxClient) RunCmd(args ...string) error {
	argList := []string{}
	if t.SocketPath != "" {
		argList = append(argList, "-S", t.SocketPath, args...)
	}
	return RunCmd("tmux", argList...)

}
