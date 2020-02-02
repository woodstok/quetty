package quetty

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
)

func GetFzfArgs() []string {
	headerArg := "--header="
	bindingArg := "--bind="
	keyMap := make(map[string]string)
	keyMap["ctrl-w"] = "word"
	keyMap["ctrl-e"] = "nospace"
	keyMap["ctrl-h"] = "hash"
	keyMap["ctrl-n"] = "num"
	keyMap["ctrl-p"] = "path"
	keyMap["ctrl-i"] = "ip"
	keyMap["ctrl-t"] = "time"

	for key, tokenType := range keyMap {
		headerArg += fmt.Sprintf("%s:%s ", key, tokenType)
		bindingArg += fmt.Sprintf("%s:reload(quetty -%s),", key, tokenType)
	}
	bindingArg = strings.TrimSuffix(bindingArg, ",")
	return []string{"--print0", bindingArg, headerArg}
}

type Printer struct {
	options *Options
}

func NewPrinter(opts *Options) (*Printer, error) {
	return &Printer{options: opts}, nil
}

func setupFzfTmuxCommand() (io.WriteCloser, *exec.Cmd) {

	var fzfCmd *exec.Cmd

	fzfCmd = exec.Command("fzf-tmux", GetFzfArgs()...)
	writer, err := fzfCmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	return writer, fzfCmd
}
func runFzfTmuxCommand(wg *sync.WaitGroup, fzfCmd *exec.Cmd) {
	wg.Add(1)
	defer wg.Done()
	tmuxClient := NewTmuxClient()
	curPane, err := tmuxClient.CurrentPane()
	if err != nil {
		log.Fatal(err)
	}
	output, err := fzfCmd.CombinedOutput()
	if err == nil && curPane != "" {
		err = tmuxClient.SendString(curPane, output)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// if init flag is passed, we need to start fzf in which case,
// all the tokens should be piped to the fzf process.
// otherwise, the expectation is to print it to stdout.
// The fzf process will reread stdout of a quetty process when it reloads
// because of a keybinding
func (p *Printer) Print(tokens Tokens) {
	var writer io.WriteCloser
	var wg sync.WaitGroup
	if p.options.init {
		fzfInputPipe, fzfCmd := setupFzfTmuxCommand()
		writer = fzfInputPipe
		go runFzfTmuxCommand(&wg, fzfCmd)
	} else {
		writer = os.Stdout
	}

	for token, _ := range tokens {
		fmt.Fprintln(writer, token)
	}
	writer.Close()
	wg.Wait()
}
