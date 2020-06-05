FROM golang:alpine
RUN mkdir /app
ADD . /app/
WORKDIR /app

RUN go mod download && go mod verify && go mod tidy && go mod vendor
RUN go build -mod=vendor -o ./bin/billing ./cmd/billing
RUN go build -mod=vendor -o ./bin/queue ./cmd/queue

RUN apk add supervisor

COPY ./.docker/supervisord/conf.d/supervisord.conf /etc/supervisord/conf.d/supervisord.conf

ENTRYPOINT ["supervisord", "--nodaemon", "--configuration", "/etc/supervisord/conf.d/supervisord.conf"]