package controllers

import (
	"eks-fargate-logger/helper"
	"time"
)

var logger = helper.Logger()

type Kubernetes struct {
	PodName        string            `json:"pod_name"`
	NamespaceName  string            `json:"namespace_name"`
	PodID          string            `json:"pod_id"`
	Labels         map[string]string `json:"labels"`
	Annotations    map[string]string `json:"annotations"`
	HostName       string            `json:"host"`
	ContainerName  string            `json:"container_name"`
	DockerID       string            `json:"docker_id"`
	ContainerHash  string            `json:"container_hash"`
	ContainerImage string            `json:"container_image"`
}

type LogRecord struct {
	Log        interface{} `json:"log"`
	Kubernetes Kubernetes  `json:"kubernetes"`
}

type Record struct {
	Data        string `json:"data"`
	Kuberenetes struct {
		PodName       string            `json:"pod_name"`
		Namespace     string            `json:"namespace"`
		PodID         string            `json:"pod_id"`
		Labels        map[string]string `json:"labels"`
		Annotations   map[string]string `json:"annotations"`
		HostName      string            `json:"host"`
		ContainerName string            `json:"container_name"`
		DockerImage   string            `json:"docker_image"`
		Timestamp     string            `json:"time"`
	} `json:"kubernetes"`
}

type Response struct {
	RequestID string   `json:"requestId"`
	Timestamp int64    `json:"timestamp"`
	Records   []Record `json:"records"`
}

type DocRecord struct {
	Stream     string         `json:"stream"`
	Docker     DockerInfo     `json:"docker,omitempty"`
	Kubernetes KubernetesInfo `json:"kubernetes,omitempty"`
	Message    string         `json:"message"`
	Timestamp  time.Time      `json:"@timestamp"`
	Tag        string         `json:"tag"`
}

type DockerInfo struct {
	ContainerID string `json:"container_id,omitempty"`
}

type KubernetesInfo struct {
	ContainerName    string            `json:"container_name,omitempty"`
	NamespaceName    string            `json:"namespace_name,omitempty"`
	PodName          string            `json:"pod_name,omitempty"`
	ContainerImage   string            `json:"container_image,omitempty"`
	ContainerImageID string            `json:"container_image_id,omitempty"`
	PodID            string            `json:"pod_id,omitempty"`
	PodIP            string            `json:"pod_ip,omitempty"`
	Host             string            `json:"host,omitempty"`
	Labels           map[string]string `json:"labels,omitempty"`
	MasterURL        string            `json:"master_url,omitempty"`
	NamespaceID      string            `json:"namespace_id,omitempty"`
	NamespaceLabels  map[string]string `json:"namespace_labels,omitempty"`
}
