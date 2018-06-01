FROM golang:1.10 as build
WORKDIR /run-occasionally
COPY main.go .
RUN go build

FROM alpine
COPY --from=build /run-occasionally/run-occasionally /usr/local/bin/
ENTRYPOINT ["run-occasionally"]
