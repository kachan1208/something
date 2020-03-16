FROM golang:1.14 as build

WORKDIR /

COPY . /

RUN go mod download

RUN GOOS=linux GOARCH=amd64 make build

FROM alpine:3.10

COPY --from=build /service ./

ENTRYPOINT ["/service"]