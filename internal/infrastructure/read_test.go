//go:build integration

package infrastructure

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"homework-1/internal/service"
	"homework-1/internal/storage"
	"io"
	"log"
	"os"
	"strings"
	"testing"
)

func TestSimpleInput(t *testing.T) {
	st, err := storage.New()
	if err != nil {
		log.Fatal(err)
	}
	defer st.CloseStorage()
	serv := service.New(st)
	s := NewSync()
	signals := make(chan os.Signal)
	oldStdout := os.Stdout
	readOut, writeOut, _ := os.Pipe()
	os.Stdout = writeOut
	input := `help
exit`
	in := bufio.NewReader(strings.NewReader(input))

	go Process(context.Background(), in, signals, s)
	s.Wg.Add(1)
	go ControlRoutines(serv, s)
	WaitForExitSignal(signals)
	WaitForExit(s)
	writeOut.Close()
	var output bytes.Buffer
	io.Copy(&output, readOut)
	os.Stdout = oldStdout

	assert.Contains(t, output.String(), "This is a command line tool")
}

func TestIncorrectInput(t *testing.T) {
	st, err := storage.New()
	if err != nil {
		log.Fatal(err)
	}
	defer st.CloseStorage()
	serv := service.New(st)
	s := NewSync()
	signals := make(chan os.Signal)
	oldStdout := os.Stdout
	_, writeOut, _ := os.Pipe()
	os.Stdout = writeOut
	var logOutput bytes.Buffer
	log.SetOutput(&logOutput)
	defer log.SetOutput(os.Stderr)
	input := `
abracadabra
accept
back
return
list
exit`
	in := bufio.NewReader(strings.NewReader(input))

	go Process(context.Background(), in, signals, s)
	s.Wg.Add(1)
	go ControlRoutines(serv, s)
	WaitForExitSignal(signals)
	WaitForExit(s)
	writeOut.Close()
	os.Stdout = oldStdout
	fmt.Println(logOutput.String())
	assert.Contains(t, logOutput.String(), "Command  doesn't exist. Try command help")
	assert.Contains(t, logOutput.String(), "date is missing")
	assert.Contains(t, logOutput.String(), "order id is missing")
	assert.Contains(t, logOutput.String(), "client id is missing")
}

func TestCorrectLogicInput(t *testing.T) {
	st, err := storage.New()
	if err != nil {
		log.Fatal(err)
	}
	defer st.CloseStorage()
	serv := service.New(st)
	s := NewSync()
	signals := make(chan os.Signal)
	oldStdout := os.Stdout
	readOut, writeOut, _ := os.Pipe()
	os.Stdout = writeOut
	var logOutput bytes.Buffer
	log.SetOutput(&logOutput)
	defer log.SetOutput(os.Stderr)
	input := `change -num_routines=4
pickup -list=1,2
exit`
	in := bufio.NewReader(strings.NewReader(input))

	go Process(context.Background(), in, signals, s)
	s.Wg.Add(1)
	go ControlRoutines(serv, s)
	WaitForExitSignal(signals)
	WaitForExit(s)
	writeOut.Close()
	var output bytes.Buffer
	io.Copy(&output, readOut)
	os.Stdout = oldStdout

	assert.Contains(t, output.String(), "changed number routines from 5 to 4")
	assert.Contains(t, logOutput.String(), "some of orders were not found")
}
