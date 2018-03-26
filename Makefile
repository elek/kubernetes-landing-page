#################################################
GOOS	:= $(shell go env GOOS)
GOARCH	:= $(shell go env GOARCH)
GOFILES	:= $(shell ls *.go |grep -v test)
GOBUILD	:= GOOS=$(GOOS) GOARCH=$(GOARCH) go build


#################################################
default:	deps test install
docker:	docker


bindata:
	go-bindata -o ./handlers/bindata.go -pkg handlers views/...

deps:
	glide install

build:
	$(GOBUILD) $(GOFILES)

install:
	GOOS=linux GOARCH=amd64 go install

docker:
	docker build -t elek/kubernetes-landing-page .

test:
	go test -v
