FROM --platform=$TARGETPLATFORM golang:alpine AS builder

ENV CGO_ENABLED=0
ENV GOOS=linux

WORKDIR /build
COPY go.mod .
COPY go.sum .
RUN go mod tidy
COPY . .
RUN CGO_ENABLED=$CGO_ENABLED go build -ldflags "-extldflags='-static'" -o base64

FROM alpine:3.23.2
WORKDIR /app
COPY --from=builder /build/base64 /app/base64
COPY --from=builder /build/templates /app/templates
COPY --from=builder /build/static /app/static
ENTRYPOINT ["/app/base64"]
