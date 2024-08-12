FROM golang:1.22-alpine AS builder

WORKDIR /app

RUN apk update && apk upgrade && \
    apk add --no-cache bash git make

# dependencies
COPY ["go.mod", "go.sum", "./"]
RUN go mod download -x

# build
COPY . ./
RUN go build -o ./bin/myApp cmd/app/main.go

FROM alpine

COPY --from=builder /app/bin/myApp /app/bin/myApp

RUN chmod +x /app/bin/myApp

CMD ["/app/bin/myApp"]