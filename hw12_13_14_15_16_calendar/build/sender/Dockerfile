# Собираем в гошке
FROM golang:1.23 as builder

ENV BIN_FILE /opt/sender/sender-app
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
        -o ${BIN_FILE} cmd/sender/*

# На выходе тонкий образ
FROM alpine:3.9

LABEL ORGANIZATION="DimRyb"
LABEL SERVICE="sender"
LABEL MAINTAINERS="dimryb@bk.ru"

ENV BIN_FILE /opt/sender/sender-app
COPY --from=builder ${BIN_FILE} ${BIN_FILE}

ENV CONFIG_FILE /etc/sender/config.yaml
COPY ./configs/sender.yaml ${CONFIG_FILE}

COPY cmd/scripts/wait-for-rabbit.sh /wait-for-rabbit.sh
RUN chmod +x /wait-for-rabbit.sh

CMD ["sh", "-c", "/wait-for-rabbit.sh && ${BIN_FILE} -config ${CONFIG_FILE}"]