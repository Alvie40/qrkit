FROM golang:1.24.3-alpine

WORKDIR /app
COPY . .
RUN go build -o server ./cmd/server

EXPOSE 8080
CMD ["./server"]
