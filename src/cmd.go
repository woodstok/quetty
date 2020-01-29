package quetty

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
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
		bindingArg += fmt.Sprintf("%s:reload(quetty -r -%s),", key, tokenType)
	}
	bindingArg = strings.TrimSuffix(bindingArg, ",")
	bindingArg = "--bind=ctrl-w:reload(ps -elf)"
	return []string{bindingArg}
}

func Run(opts *Options) {
	f, err := os.OpenFile("quetty.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	tmuxClient := NewTmuxClient()
	curWindow, err := tmuxClient.CurrentWindow()
	if err != nil {
		log.Fatal(err)
	}

	paneList, err := tmuxClient.ListWindowPanes(curWindow)
	if err != nil {
		log.Fatal(err)
	}

	pipeList := []io.Reader{}
	for _, paneId := range paneList {
		panePipeR, panePipeW := io.Pipe()
		pipeList = append(pipeList, panePipeR)
		go tmuxClient.WriteFromPane(paneId, panePipeW)
	}
	combinedReader, combinedWriter := io.Pipe()

	go func() {
		io.Copy(combinedWriter, io.MultiReader(pipeList...))
		combinedWriter.Close()
	}()

	tokenMgr, err := NewTokenMgr(opts)
	if err != nil {
		log.Fatal(err)
	}

	tokens, err := tokenMgr.Process(combinedReader)
	if err != nil {
		log.Fatal(err)
	}

	printer, err := NewPrinter(opts)
	if err != nil {
		log.Fatal(err)
	}

	fzfCmd := exec.Command("fzf-tmux", GetFzfArgs()...)
	fzfInputPipe, err := fzfCmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	go printer.Print(fzfInputPipe, tokens)
	output, err := fzfCmd.CombinedOutput()
	if err != nil {
		log.Fatal(output, err)
	}
	if err != nil {
		log.Fatal(err)
	}
}
