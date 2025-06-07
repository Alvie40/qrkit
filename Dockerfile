FROM golang:1.23

WORKDIR /app
COPY . .
RUN go build -o out ./cmd/server

EXPOSE 8080
CMD ["./out"]
