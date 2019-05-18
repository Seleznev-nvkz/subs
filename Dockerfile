# stage 1
FROM golang:1.12-alpine AS build_base

RUN apk add --no-cache ca-certificates git
WORKDIR /src

COPY ./go.mod ./go.sum ./
RUN go mod download
COPY ./ ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags="-s -w" -installsuffix cgo -o /app .

# stage 2
FROM scratch as subtitles_helper
WORKDIR /root/
COPY --from=build_base /app /app
COPY ./static ./static

ENTRYPOINT ["/app"]
