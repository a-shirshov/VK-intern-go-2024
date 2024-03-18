server:
	go build -o server -v ./main.go

swagger:
	GO111MODULE=off swagger generate spec -o ./swagger.yaml --scan-models

mock:
	./mockgen.sh

cover:
	go test -cover -coverprofile=cover.out -coverpkg=./... ./...
	cat cover.out | fgrep -v "main.go" | fgrep -v "mock.go" | fgrep -v "swagger_models.go" | fgrep -v "postgres.go" > cover1.out
	go tool cover -func=cover1.out