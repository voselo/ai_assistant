# Стадия сборки
FROM golang:1.22-alpine as builder

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh make && \
    apk add --no-cache tzdata

# Копируем локальный код в образ контейнера.
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download -x

# Копируем весь проект в контейнер
COPY . .

# Собираем приложение
RUN go build -o /app/main ./cmd/app/main.go

# Финальная стадия
FROM alpine

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем собранное приложение из стадии сборки
COPY --from=builder  /app/main .

# Указываем команду по умолчанию для запуска приложения
CMD ["/app/main"]