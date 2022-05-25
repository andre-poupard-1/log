package buffer

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

type PostBackBuffer[T any] struct {
	Incoming *chan T
	Interval *chan int
	Buffer []T
	BufferLimit int
	PostbackUrl string
}

func (buf *PostBackBuffer[T]) waitForEvents() {
	for {
		select {
			case event := <-*buf.Incoming:
				buf.HandleIncomingEvent(event)
			case <-*buf.Interval:
				buf.HandleIntervalEvent()
		}
	}
}

func (buf *PostBackBuffer[T]) startIntervalTicker () {
	intervalTicker := time.NewTicker(2 * time.Second)
	go func() {
		for {
			<-intervalTicker.C
			*buf.Interval <- 1
		}
	}()
}

func (buf *PostBackBuffer[T]) HandleIntervalEvent() {
	if len(buf.Buffer) != 0 {
		buf.EmptyBuffer()
	}
}

func (buf *PostBackBuffer[T]) EmptyBuffer() {
	oldBuf := buf.ResetBuffer()
	go buf.SendPostBack(oldBuf)
}

func (buf *PostBackBuffer[T]) HandleIncomingEvent(event T) {
	buf.Buffer = append(buf.Buffer, event)
	if buf.AtCapacity() {
		buf.EmptyBuffer()
	}
}

func (buf *PostBackBuffer[T]) AtCapacity() bool {
	return len(buf.Buffer) == buf.BufferLimit
}

func (buf *PostBackBuffer[T]) SendPostBack(body []T) (*http.Response, error) {
	postBody, _ := json.Marshal(body)
	resp, err := http.Post(buf.PostbackUrl, "application/json", bytes.NewBuffer(postBody))
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (buf *PostBackBuffer[T]) AttemptSendPostBack(body []T) {
	attempts := 0
	const MAX_ATTEMPTS = 3
	for {
		_, err := buf.SendPostBack(body)
		if err != nil {
			if attempts == MAX_ATTEMPTS {
				break
			}
		} else {
			break
		}

		attempts += 1
		time.Sleep(time.Second * 2)
	}
}

func (buf *PostBackBuffer[T]) ResetBuffer() ([]T) {
	cpy := make([]T, len(buf.Buffer))
	copy(cpy, buf.Buffer)
	buf.Buffer = make([]T, 0, buf.BufferLimit)
	return cpy
}