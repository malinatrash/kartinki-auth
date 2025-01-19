FROM golang:1.22.4-alpine

WORKDIR /src
# Copy proto files
COPY proto proto/

WORKDIR /src/app
# Copy auth files
COPY auth/go.mod auth/go.sum ./
RUN go mod download
COPY auth/ .
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go
CMD ["./main"]
