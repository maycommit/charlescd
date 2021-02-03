FROM golang:1.14-alpine as builder

# RUN apk update && apk add gcc libc-dev git

WORKDIR /app/circlerr

COPY ./go.mod .
COPY ./go.sum .

RUN go mod download

COPY . .

RUN go build -o ./out/circlerr cmd/controller/*.go

# -----------------------------------------------

FROM alpine:3.12.0

COPY --from=builder /app/circlerr/out/circlerr .

EXPOSE 8080

ENTRYPOINT ["./circlerr"]