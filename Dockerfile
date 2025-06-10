FROM golang:1.23

WORKDIR /app

# Copie os arquivos de dependências primeiro para cache eficiente
COPY go.mod .
COPY go.sum .

RUN go mod download

# Agora copie o resto do código
COPY . .

RUN go build -o out ./cmd/server

EXPOSE 8080
CMD ["./out"]
