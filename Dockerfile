FROM golang:alpine as builder
RUN mkdir /build 
ADD . /build/
WORKDIR /build 
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-w -s -extldflags "-static"' -o main .

FROM scratch
COPY --from=builder /build/main /app/shortflake
WORKDIR /app
ENTRYPOINT ["./shorflake"]