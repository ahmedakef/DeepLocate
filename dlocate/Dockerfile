FROM golang:latest

WORKDIR /app

COPY ./ /app

RUN go mod download

RUN go get github.com/githubnemo/CompileDaemon

# RUN go build

ENTRYPOINT CompileDaemon
#ENTRYPOINT ./dlocate -o index -d /home/