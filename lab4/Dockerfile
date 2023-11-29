# Используем официальный образ Golang 1.21.3
FROM golang:1.21.3

# Устанавливаем рабочую директорию
WORKDIR /

# Копируем файлы проекта в рабочую директорию
COPY . .

# Запускаем go mod tidy
RUN go mod tidy

# Запускаем сервер
CMD ["go", "run", "cmd/server_2/server.go"]