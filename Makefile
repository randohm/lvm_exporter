all: lvm_exporter

lvm_exporter: main.go
	GOOS=linux GOARCH=amd64 go build
