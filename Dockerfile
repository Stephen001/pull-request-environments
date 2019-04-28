FROM golang:alpine as builder
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN apk add --no-cache git
RUN go get gopkg.in/go-playground/webhooks.v5/github
RUN go get github.com/fsnotify/fsnotify
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main .
FROM scratch
COPY --from=builder /build/main /app/
WORKDIR /app
CMD ["./main"]
EXPOSE 3000
