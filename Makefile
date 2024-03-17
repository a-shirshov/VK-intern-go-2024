server:
	go build -o server -v ./main.go

swagger:
	GO111MODULE=off swagger generate spec -o ./swagger.yaml --scan-models