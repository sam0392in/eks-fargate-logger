package tests

import (
	"eks-fargate-logger/models"
	"github.com/gin-gonic/gin"
	"time"
)

func EsSend(c *gin.Context) {
	logger.Info("Serving Path: ", c.FullPath())
	// Create a time object with the desired date and time
	t := time.Now()

	// Format the time object as a string in the desired format
	formattedTime := t.Format("2006-01-02T15:04:05.999999999-07:00")

	var data = map[string]interface{}{
		"stream": "stdout",
		"docker": map[string]interface{}{
			"container_id": "fd716ecc7a7cec069776e01f1c9111148a4ed6b1dba3aac71774e7841e74a538",
		},
		"kubernetes": map[string]interface{}{
			"container_name":     "samapp", // Fill in the missing container_name field
			"namespace_name":     "cluster-services",
			"pod_name":           "samfargateapp-5f4bf69998-4bvcw",
			"container_image":    "123445566778.dkr.ecr.eu-west-1.amazonaws.com/sam:test-go_0.6",
			"container_image_id": "123445566778.dkr.ecr.eu-west-1.amazonaws.com/sam@sha256:1305dff46ee2d73c162872bf1fae5116aab69d43212d84d15fa2ae018acf7fa6",
			"host":               "fargate-ip-172-23-3-181.eu-west-1.compute.internal",
			"labels": map[string]interface{}{
				"app":                               "samapp",
				"eks.amazonaws.com/fargate-profile": "devops",
				"nature":                            "serverless",
				"pod-template-hash":                 "5f4bf69998",
			},
			// Fill in the missing namespace_id and namespace_labels fields
			"namespace_id": "3c913aa7-5e01-46dc-a671-1ea1dae2939e",
			"namespace_labels": map[string]interface{}{
				"kubernetes.io/metadata.name": "cluster-services",
			},
		},
		// Fill in the missing fields
		"logger_name": "eks-fargate-logger",
		"namespace_labels": map[string]interface{}{
			"kubernetes.io/metadata.name": "cluster-services",
		},
		"@timestamp": formattedTime,
	}

	index := "k8s-cluster-services--2024.03.26"
	//models.BulkIndexing(index, data)
	formattedData := models.FormatData(index, data)
	dataChan <- formattedData // Send data to buffer channel
}
