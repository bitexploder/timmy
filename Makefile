GOOS=linux
GOARCH=amd64
MITMHOST=mitm
DEPLOY_PATH=/home/debian/timmy/timmy
SSH_ARGS=-F ~/.ssh/mitm_config
FILES=timmy.go orig_dest.go
OPTS=-o timmy

linux:
	env GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(OPTS) $(FILES)

linux32: GOARCH=386
linux32: 	
	env GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(OPTS) $(FILES)

deploy: linux
	ssh $(SSH_ARGS) mitm "sudo supervisorctl stop timmy"
	scp $(SSH_ARGS) timmy mitm:$(DEPLOY_PATH)
	ssh $(SSH_ARGS) mitm "sudo supervisorctl start timmy"

local:
	go install
