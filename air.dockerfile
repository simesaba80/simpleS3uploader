FROM golang:1.25.5-alpine3.23
WORKDIR /app

RUN go install github.com/air-verse/air@latest
CMD ["air"]