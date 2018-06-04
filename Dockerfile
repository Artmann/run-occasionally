FROM golang:1.10 as build

WORKDIR /go/src/app
COPY main.go .

ENV CGO_ENABLED=0
RUN go get -d -v ./...
RUN go install -v ./...

FROM alpine
COPY --from=build /go/bin/app /usr/bin/run-occasionally
ENTRYPOINT ["run-occasionally"]
