SERVER_NAME=thingful_hackweek
CLI_NAME=thingful_hackweek

build : linux windows 

linux: tmp/build/$(SERVER_NAME)-linux-amd64 tmp/build/$(CLI_NAME)-linux-amd64
windows: tmp/build/$(SERVER_NAME)-windows-amd64.exe tmp/build/$(CLI_NAME)-windows-amd64.exe

tmp/build/$(SERVER_NAME)-linux-amd64:
	GOOS=linux GOARCH=amd64 go build  -o $(@) ./cmd/server/

tmp/build/$(SERVER_NAME)-windows-amd64.exe:
	GOOS=windows GOARCH=amd64  go build -o $(@) ./cmd/server/


