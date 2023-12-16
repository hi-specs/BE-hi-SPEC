FROM golang:1.21.4-alphine

COPY . /app

WORKDIR /app

RUN go mod tidy

RUN go build -o server .

CMD [ "/app/server" ]