FROM golang:1.16 as builder

RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN go build -o main .

FROM gradle:6.8.3-jdk8

COPY --from=builder /build/main /app/

ENV PORT=8080
EXPOSE 8080

WORKDIR /app
ENTRYPOINT ["./main"]
