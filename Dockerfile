# ---------- build stage ----------
    FROM golang:1.22-alpine AS builder

    # Устанавливаем системные зависимости
    RUN apk add --no-cache git
    
    WORKDIR /app
    
    # Копируем go.mod и go.sum отдельно (для кеша)
    COPY go.mod go.sum ./
    RUN go mod download
    
    # Копируем исходники
    COPY . .
    
    # Собираем бинарник
    RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bot
    
    
    # ---------- runtime stage ----------
    FROM alpine:3.19
    
    # UTF-8 (важно для русских имён файлов)
    ENV LANG=C.UTF-8
    ENV LC_ALL=C.UTF-8
    
    # Устанавливаем runtime-зависимости
    RUN apk add --no-cache \
        ffmpeg \
        python3 \
        py3-pip \
        ca-certificates
    
    # Устанавливаем yt-dlp
    RUN pip install --no-cache-dir yt-dlp
    
    WORKDIR /app
    
    # Копируем бинарник из build stage
    COPY --from=builder /app/bot /app/bot
    
    # Render НЕ требует EXPOSE, но пусть будет
    EXPOSE 8080
    
    # Запуск бота
    CMD ["./bot"]
    