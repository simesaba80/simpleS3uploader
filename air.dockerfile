FROM golang:1.25.5-alpine3.23
WORKDIR /app
RUN apk --update add tzdata && \
    cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime && \
    apk del tzdata && \
    rm -rf /var/cache/apk/*
RUN go install github.com/air-verse/air@v1.62.0
CMD ["air"]