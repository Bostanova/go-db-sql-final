FROM golang:1.23

WORKDIR /app

COPY . .

#go mod download загружает все модули и пакеты, указанные в go.mod и
#go.sum (используем вместо go mod init)
RUN go mod download && go mod verify

# CGO_ENABLED=0 отключаем поддержку библиотек языка С, т.к. в этом проекте они не нужны
# GOOS=linux указываем, что проект собирается для линукс
# GOARCH=amd64 указываем платформу, для которой будет компилироваться исполняемый файл
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

RUN go build -o /main ./...

CMD ["/main"]