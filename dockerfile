FROM golang:1.21-alpine as server_build
RUN apk add --no-cache \
    make
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN make server
EXPOSE 8080
CMD ["./server"]