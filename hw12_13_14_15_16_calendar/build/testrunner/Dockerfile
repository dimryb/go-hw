FROM golang:1.23 as builder

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go test -c -o /integration-tests ./tests/integration

FROM alpine:3.9

LABEL ORGANIZATION="DimRyb"
LABEL SERVICE="testrunner"
LABEL MAINTAINERS="dimryb@bk.ru"

WORKDIR /app

COPY --from=builder /integration-tests /app/integration-tests
COPY cmd/scripts/wait-for-rabbit.sh /wait-for-rabbit.sh
RUN chmod +x /wait-for-rabbit.sh

ENTRYPOINT ["sh", "-c", "/wait-for-rabbit.sh && /app/integration-tests"]
