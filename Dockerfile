FROM golang:alpine as go
WORKDIR /app
ENV GO111MODULE=on

COPY go.mod .
RUN go mod download

COPY . .
RUN go build -o kurushimi ./internal/app/kurushimi

FROM alpine

WORKDIR /app

COPY --from=go /app/kurushimi .
COPY --from=go /app/internal/config/dynamic/config.yaml .
CMD ["./kurushimi"]