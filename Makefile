ALL:
	GOOS=linux GOARCH=amd64 go build -o ./bin/csi-demo-driver ./cmd
