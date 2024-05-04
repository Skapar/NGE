FROM golang:1.21.6 AS build

WORKDIR /src

COPY . .

RUN go mod download

RUN CGO_ENABLED = 0 GOOS=linux GOARCH=amd64 go build -a -installsufix cgo -a app .

FROM ubuntu:latest

RUN apt-get update
RUN apt-get install ca-certificates -y
RUN update-ca-certificates

WORKDIR /app

COPY --from=build /src/cmd/nge/ .

EXPOSE 8090

CMD ["./app"]

