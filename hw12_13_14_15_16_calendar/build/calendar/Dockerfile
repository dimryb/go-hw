FROM golang:1.23 as builder
WORKDIR /app
COPY ../.. .
RUN make build

FROM gcr.io/distroless/static-debian12
WORKDIR /app
COPY --from=builder /app/bin/calendar /app/calendar
COPY ../../configs/calendar.yaml /configs/

CMD ["/app/calendar", "--config=/configs/calendar.yaml"]