FROM golang:1.21.1

WORKDIR /app

RUN go install github.com/cosmtrek/air@latest

COPY go.mod .
COPY go.sum .
RUN go mod download


COPY . .

CMD ["air"]