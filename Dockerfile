FROM golang:alpine as go
WORKDIR /app
ENV GO111MODULE=on

COPY go.mod .
RUN go mod download

COPY . .
RUN go build -o kurushimi ./cmd/kurushimi

FROM alpine

WORKDIR /app

COPY ./run/config.yaml ./config.yaml
COPY --from=go /app/kurushimi ./kurushimi
CMD ["./kurushimi"]