package middleware

import (
	"eks-fargate-logger/helper"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

var logger = helper.Logger()

func Firewall(c *gin.Context) {
	// Check for the presence of the X-Amz-Firehose-Source-Arn header in every request
	if forwardedFor := c.Request.Header.Get("X-Amz-Firehose-Source-Arn"); forwardedFor == "" {
		logger.WithFields(logrus.Fields{
			"message": "Request without X-Amz-Firehose-Source-Arn header dropped",
		}).Warn("Unauthorized request")
		c.AbortWithStatus(http.StatusForbidden) // Respond with 403 Forbidden if request don't have desired header.
		return
	}

	path := c.Request.URL.Path
	method := c.Request.Method

	// Implement custom logging rather than using Default GIN Logger
	logger.WithFields(logrus.Fields{
		"path":   path,
		"method": method,
	}).Info("Incoming request")

	c.Next()
}
