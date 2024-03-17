FROM golang:alpine

WORKDIR /app
COPY . /app
RUN go build -o flow

ENTRYPOINT ["./flow"]
