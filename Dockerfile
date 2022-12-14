# Build the go application into a binary
FROM golang:1.19.4-alpine as builder
WORKDIR /app
COPY . ./
RUN go build -o dongato .

FROM alpine
COPY --from=builder /app/dongato .
ENTRYPOINT ["/dongato"]
