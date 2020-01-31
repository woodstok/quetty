package quetty

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

const bufferFilePath = "/tmp/quetty-buffer"

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
	return []string{"--print0", bindingArg, headerArg}
}

var logger *log.Logger

func Run(opts *Options) {
	f, err := os.OpenFile("/tmp/quetty.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	logger = log.New(f, "prefix", log.LstdFlags)
	logger.Printf("Starting quetty")
	defer f.Close()

	tmuxClient := NewTmuxClient()
	var curPane string
	var inputStream io.ReadCloser

	if !opts.reload {
		curWindow, err := tmuxClient.CurrentWindow()
		if err != nil {
			log.Fatal(err)
		}

		curPane, err = tmuxClient.CurrentPane()
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

		os.Remove(bufferFilePath)
		tmpFile, err := os.OpenFile(bufferFilePath, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}

		go func() {
			io.Copy(io.MultiWriter(combinedWriter, tmpFile), io.MultiReader(pipeList...))
			combinedWriter.Close()
			tmpFile.Close()
		}()
		inputStream = combinedReader

	} else {
		inputReader, err := os.OpenFile(bufferFilePath, os.O_RDONLY, 0644)
		inputStream = inputReader
		if err != nil {
			log.Fatal(err)
		}

	}

	logger.Printf("running with opts %+v", opts)
	tokenMgr, err := NewTokenMgr(opts)
	if err != nil {
		log.Fatal(err)
	}

	tokens, err := tokenMgr.Process(inputStream)
	if err != nil {
		log.Fatal(err)
	}
	inputStream.Close()

	printer, err := NewPrinter(opts)
	if err != nil {
		log.Fatal(err)
	}

	var fzfCmd *exec.Cmd
	if !opts.reload {

		fzfCmd = exec.Command("fzf-tmux", GetFzfArgs()...)
		fzfInputPipe, err := fzfCmd.StdinPipe()
		if err != nil {
			log.Fatal(err)
		}
		go printer.Print(fzfInputPipe, tokens)
		output, err := fzfCmd.CombinedOutput()
		if err == nil && curPane != "" {
			err = tmuxClient.SendString(curPane, output)
			if err != nil {
				log.Fatal(err)
			}
		}

	} else {
		printer.Print(os.Stdout, tokens)
	}
	os.Exit(0)
}
