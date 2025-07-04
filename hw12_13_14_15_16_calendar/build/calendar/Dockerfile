# Собираем в гошке
FROM golang:1.23 as builder

ENV BIN_FILE /opt/calendar/calendar-app
ENV CODE_DIR /go/src/

WORKDIR ${CODE_DIR}

# Кэшируем слои с модулями
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

# Собираем статический бинарник Go (без зависимостей на Си API),
# иначе он не будет работать в alpine образе.
ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE} cmd/calendar/*

# На выходе тонкий образ
FROM alpine:3.9

LABEL ORGANIZATION="DimRyb"
LABEL SERVICE="calendar"
LABEL MAINTAINERS="dimryb@bk.ru"

ENV BIN_FILE /opt/calendar/calendar-app
COPY --from=builder ${BIN_FILE} ${BIN_FILE}

ENV CONFIG_FILE /etc/calendar/config.yaml
COPY ./configs/calendar.yaml ${CONFIG_FILE}

COPY migrations /app/migrations

RUN apk add --no-cache curl

CMD ${BIN_FILE} -config ${CONFIG_FILE}
