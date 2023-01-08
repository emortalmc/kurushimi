FROM golang:alpine as go
WORKDIR /app
ENV GO111MODULE=on

COPY go.mod .
RUN go mod download

COPY . .
RUN go build -o kurushimi ./internal/director

FROM alpine

COPY --from=go /app/kurushimi /app/kurushimi
CMD ["/app/kurushimi"]