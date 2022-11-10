# Собираем в гошке
FROM golang:1.16.2 as build

ENV BIN_FILE /opt/calendar/calendar-app
ENV CODE_DIR /go/src/

WORKDIR ${CODE_DIR}

# БД (скрипт по автогенерации UUID)
#COPY /internal/utils/db_init_scripts.sql /docker-entrypoint-initdb.d

# Добавляем инструменты миграции
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Добавляем протобаф + grpc
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

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

LABEL ORGANIZATION="OTUS Online Education"
LABEL SERVICE="calendar"
LABEL MAINTAINERS="student@otus.ru"

ENV BIN_FILE "/opt/calendar/calendar-app"
COPY --from=build ${BIN_FILE} ${BIN_FILE}

ENV CONFIG_FILE /etc/calendar/config.toml
COPY ./configs/config.toml ${CONFIG_FILE}

CMD ${BIN_FILE} -config ${CONFIG_FILE}
