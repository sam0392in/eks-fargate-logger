.SILENT: build
.PHONY:	release build install uninstall

APPNAME = eks-fargate-logger
IMAGE = eks-fargate-logger
NAMESPACE = cluster-services
PLATFORM = amd64
TAG ?= dev-tag

# Default to dev environment
ENV ?= dev
ifeq ($(ENV),dev)
	TAG = dev-tag
	PREFIX = 123445566778.dkr.ecr.eu-west-1.amazonaws.com
	VALUEFILE = environments/values-dev.yaml
else
	TAG = live-tag
	PREFIX = 998877665544.dkr.ecr.eu-west-1.amazonaws.com
	VALUEFILE = environments/values-live.yaml
endif

deploy: build push release

build:
	docker build --platform=linux/$(PLATFORM) -t $(PREFIX)/$(IMAGE):$(TAG) .

push:
	docker push $(PREFIX)/$(IMAGE):$(TAG)

release:
	helm upgrade -i $(APPNAME) -n $(NAMESPACE) ./k8s -f ./k8s/${VALUEFILE} \
	--set-string image.tag=$(TAG)

dryrun:
	helm upgrade -i $(APPNAME) -n $(NAMESPACE) ./k8s -f ./k8s/${VALUEFILE} \
	--set-string image.tag=$(TAG) \
	--dry-run

uninstall:
	helm uninstall $(APPNAME) -n $(NAMESPACE)
