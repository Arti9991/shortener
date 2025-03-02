FROM golang:1.22

WORKDIR /usr/src/shortener

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download

COPY . .
WORKDIR /usr/src/shortener/cmd/shortener
EXPOSE 8082
ENV SERVER_ADDRESS=:8082

RUN go build -v -o /usr/local/bin/shortener ./...

ENTRYPOINT ["shortener"]