# Используем официальный образ Golang 1.21.3
FROM golang:1.21.3

RUN apt-get update && apt-get install -y iputils-ping

# Устанавливаем рабочую директорию
WORKDIR /

# Копируем файлы проекта в рабочую директорию
COPY . .

# Запускаем сервер
CMD ["go", "run", "cmd/server_2/server.go"]