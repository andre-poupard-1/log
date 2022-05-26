package post

import (
	"bytes"
	"encoding/json"
	"fmt"

	"main/config"
	logger "main/middleware"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type PostBackBuffer[T any] struct {
	Incoming *chan T
	Interval *chan bool
	Buffer []T
	BufferLimit int
	PostbackUrl string
	Logger *zap.Logger
}

func (buf *PostBackBuffer[T]) waitForEvents() {
	buf.startIntervalTicker()
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
	intervalTicker := time.NewTicker(time.Duration(config.GetConfig().BatchSecondInterval) * time.Second)
	go func() {
		for {
			<-intervalTicker.C
			*buf.Interval <- true
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
	go buf.AttemptSendPostBack(oldBuf)
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
	beforeReq := time.Now()
	resp, err := http.Post(buf.PostbackUrl, "application/json", bytes.NewBuffer(postBody))	
	duration := time.Since(beforeReq)

	if resp != nil {
		buf.LogPostBack(duration.Nanoseconds(), resp.StatusCode, len(body))
		defer resp.Body.Close()
	}

	if err != nil {
		return nil, err
	}
	// handle case where error is not thrown
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("%s gave %d response status.", buf.PostbackUrl, resp.StatusCode)
	}

	return resp, err
}

func (buf *PostBackBuffer[T]) LogPostBack(duration int64, statusCode int, bufferLen int) {
	buf.Logger.Info(logger.OUTGOING_POST_REQUEST,
		zap.Int64("nanosecond_duration", duration),
		zap.Int("response_code", statusCode),
		zap.Int("buffer_size", bufferLen),
	)
}

func (buf *PostBackBuffer[T]) AttemptSendPostBack(body []T) {
	attempts := 1
	const MAX_ATTEMPTS = 3
	const SECONDS_BACKOFF = 2
	for {
		_, err := buf.SendPostBack(body)
		if err != nil {
			if attempts == MAX_ATTEMPTS {
				buf.LogFinalBackoffRequest(body)
				break
			}
		} else {
			break
		}

		attempts += 1
		time.Sleep(time.Second * SECONDS_BACKOFF)
	}
}

func (buf *PostBackBuffer[T]) ResetBuffer() ([]T) {
	cpy := make([]T, len(buf.Buffer))
	copy(cpy, buf.Buffer)
	buf.Buffer = make([]T, 0)
	return cpy
}

func (buf *PostBackBuffer[T]) LogFinalBackoffRequest(body []T) {
	stringBody, _ := json.Marshal(body)
	buf.Logger.Error(logger.FAILED_FINAL_BACKOFF_REQUEST,
		zap.String("url", buf.PostbackUrl),
		zap.String("body", string(stringBody)),
		zap.Int("buffer_size", len(body)),
	)	
}