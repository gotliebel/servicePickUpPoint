//go:build !integration

package kafka

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"sync"
	"testing"
	"time"
)

func TestProduceEventToKafka(t *testing.T) {
	producer := ProduceEventToKafka()
	oldStdout := os.Stdout
	readOut, writeOut, _ := os.Pipe()
	os.Stdout = writeOut
	event := Event{
		CreatedAt: time.Now(),
		Method:    "accept",
		RawQuery:  "-order_id=1 client_id=2 stored_until=12-12-2024",
	}
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		SendEvent(producer, event)
	}()
	go func() {
		defer wg.Done()
		SendEvent(producer, Event{CreatedAt: time.Now(), Method: "exit", RawQuery: "exited"})
	}()
	ConsumeEventsFromKafka()
	writeOut.Close()
	var output bytes.Buffer
	io.Copy(&output, readOut)
	os.Stdout = oldStdout
	wg.Wait()
	assert.Contains(t, output.String(), "Method:accept")
	assert.Contains(t, output.String(), "RawQuery:-order_id=1 client_id=2 stored_until=12-12-2024")
	assert.Contains(t, output.String(), "Method:exit")
}
