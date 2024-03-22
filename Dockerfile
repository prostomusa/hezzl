FROM golang:1.22

WORKDIR /app

COPY go.mod go.sum ./

COPY . .
RUN go mod tidy

ENTRYPOINT ["go", "run", "cmd/hezzl/main.go"]