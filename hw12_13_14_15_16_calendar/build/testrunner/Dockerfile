FROM golang:1.23 as builder

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go test -c -o /integration-tests ./tests/integration

FROM alpine:3.9

WORKDIR /app
COPY --from=builder /integration-tests /app/integration-tests

ENTRYPOINT ["/app/integration-tests"]
