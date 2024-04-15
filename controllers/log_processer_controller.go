package controllers

import (
	"eks-fargate-logger/models"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"io"
	"strings"
	"time"
)

func decodeData(data string) LogRecord {
	dt, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		logger.Error("error decoding data, ERROR: " + err.Error())
	}
	var rec LogRecord
	err = json.Unmarshal(dt, &rec)
	if err != nil {
		rootCause := errors.Unwrap(err)
		logger.Error(rootCause.Error())
	}
	return rec
}

func identifyIndex(namespace string, date time.Time) string {
	return "k8s-" + namespace + "--" + string(date.Format("2006.01.02"))
}

func extractData(data string) (string, map[string]interface{}) {
	decodedData := decodeData(data)

	logger.Debug("decoded data: ", decodedData.Log.(string))
	// decoded data = "log":"2024-03-22T10:12:23.4464632Z stdout F {\"Path\":\"/api/\",\"Response\":404,\"level\"......

	parts := strings.SplitN(decodedData.Log.(string), "F", 2)
	var logMessage string
	// Check if the split was successful (length of parts should be 2)
	if len(parts) == 2 {
		logMessage = parts[1]
	} else {
		logger.Error("String not split correctly")
	}

	var extraKeyPairs map[string]interface{}
	if err := json.Unmarshal([]byte(logMessage), &extraKeyPairs); err != nil {
		logger.Error("Error unmarshalling log message:", err)
	}

	now := time.Now()
	formattedTime := now.Format("2006-01-02T15:04:05.999999999Z")

	index := identifyIndex(decodedData.Kubernetes.NamespaceName, now)

	inputMap := map[string]interface{}{
		"@timestamp": formattedTime,
		"stream":     "stdout",
		"docker": map[string]interface{}{
			"container_id": decodedData.Kubernetes.DockerID,
		},
		"kubernetes": map[string]interface{}{
			"container_name":     decodedData.Kubernetes.ContainerName,
			"pod_name":           decodedData.Kubernetes.PodName,
			"container_image":    decodedData.Kubernetes.ContainerImage,
			"container_image_id": decodedData.Kubernetes.ContainerHash,
			"namespace_name":     decodedData.Kubernetes.NamespaceName,
			"host":               decodedData.Kubernetes.HostName,
			"labels":             decodedData.Kubernetes.Labels,
			"logger_name":        "eks-fargate-logger",
			"namespace_labels": map[string]interface{}{
				"kubernetes_io/metadata_name": decodedData.Kubernetes.NamespaceName,
			},
			"namespace_id": decodedData.Kubernetes.NamespaceName,
		},
	}
	// Add additional key-value pairs to inputMap passed from pod Logs.
	for key, value := range extraKeyPairs {
		inputMap[key] = value
	}
	return index, inputMap
}

func indexRequest(data string) {
	index, doc := extractData(data)
	dataChan <- models.FormatData(index, doc) // data formatted as per Elasticsearch input doc and sent to data buffer channel
}

func ProcessLogs(c *gin.Context) {
	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logger.Error("Error reading request body")
	}
	defer c.Request.Body.Close()

	var res Response
	err = json.Unmarshal(data, &res)
	if err != nil {
		logger.Error("Error unmarshalling data, Error: " + err.Error())
	}
	for _, r := range res.Records {
		indexRequest(r.Data)
	}
}
