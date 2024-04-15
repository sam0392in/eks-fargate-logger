package controllers

import (
	"eks-fargate-logger/models"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"time"
)

var (
	dataChan   chan esapi.IndexRequest // Channel to buffer incoming Elasticsearch index requests
	flushTimer *time.Ticker            // Timer to trigger periodic flushing of the buffer
)

func init() {
	// Create a channel with a buffer size of 50 to hold incoming index requests
	dataChan = make(chan esapi.IndexRequest, 50)

	// Create a timer that sends a signal every 5 seconds
	flushTimer = time.NewTicker(5 * time.Second)
}

func FlushBuffer() {
	var batch []esapi.IndexRequest
	for {
		select {
		// In either of the case data is flushed to ES
		case data := <-dataChan: // Receive data from the channel
			batch = append(batch, data)
			logger.Debugf("Received data in buffer")
			logger.Debugf("length of batch: " + string(rune(len(batch))))
		case <-flushTimer.C: // Receive signal from the timer
			if len(batch) > 0 {
				logger.Debugf("ReceivedInitiating flush")
				models.Flush(batch) // send the batch to Elasticsearch
				batch = nil
			}
		}
	}
}
