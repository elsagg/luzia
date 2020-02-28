## build stage
FROM golang:1.13-alpine as builder

WORKDIR /app

ENV GO111MODULE=auto

COPY . .

RUN rm .env

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

## final stage
FROM scratch

COPY --from=builder /app/luzia /app/

EXPOSE 8080

ENTRYPOINT ["/app/luzia"]