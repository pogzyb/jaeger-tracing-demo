FROM golang:1.16.6-alpine as compile

WORKDIR /app
COPY . .
RUN go mod tidy

RUN apk update && apk add git

RUN go build -o server .
RUN chmod u+x server

FROM golang:1.16.6-alpine

COPY --from=compile /app/server server

ENTRYPOINT "./server"