FROM golang:1.22

ENV GO111MODULE=on
ENV GOPROXY="https://goproxy.cn"
ENV ENV=prod

WORKDIR /app

COPY . .

RUN go install github.com/air-verse/air@v1.52.3

RUN go build -o app

EXPOSE 8080

CMD ["air", "-c", ".air.toml"]
