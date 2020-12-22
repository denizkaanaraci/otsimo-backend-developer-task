# builder image
FROM golang:1.15.6-alpine as builder
RUN mkdir /build
ADD . /build
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux go build -a -o app .

# final image
FROM alpine:3.12
COPY --from=builder /build/app .
ENTRYPOINT [ "./app" ]
