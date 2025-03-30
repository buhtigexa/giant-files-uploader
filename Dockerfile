FROM golang:1.24.0 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o giant-files-uploader main.go

FROM scratch

COPY --from=builder /app/giant-files-uploader /giant-files-uploader

ENTRYPOINT ["/giant-files-uploader"]
