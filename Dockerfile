# ---------- build stage ----------
    FROM golang:1.22-alpine AS builder

    RUN apk add --no-cache git
    
    WORKDIR /app
    
    COPY go.mod go.sum ./
    RUN go mod download
    
    COPY . .
    
    RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bot
    
    
    # ---------- runtime stage ----------
    FROM alpine:3.19
    
    ENV LANG=C.UTF-8
    ENV LC_ALL=C.UTF-8
    
    # Устанавливаем зависимости + pipx
    RUN apk add --no-cache \
        ffmpeg \
        python3 \
        py3-pip \
        pipx \
        ca-certificates
    
    # pipx требует PATH
    ENV PATH="/root/.local/bin:$PATH"
    
    # Устанавливаем yt-dlp ПРАВИЛЬНО
    RUN pipx install yt-dlp
    
    WORKDIR /app
    
    COPY --from=builder /app/bot /app/bot
    
    EXPOSE 8080
    
    CMD ["./bot"]
    