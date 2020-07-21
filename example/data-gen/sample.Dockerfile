FROM golang:1.13

LABEL maintainer="nssvlr@gmail.com"

WORKDIR /app

# COPY and download dependencies
COPY ./go.mod .
RUN go mod download


COPY . .

RUN go build -o app

CMD ["./app"]

