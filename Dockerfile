FROM golang:1.23-alpine3.20 AS base
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

FROM base AS builder
WORKDIR /src
COPY . .
ARG app_name
ENV APP_NAME $app_name
RUN CGO_ENABLED=0 go build -o "/app/service" "./services/${APP_NAME}"

FROM gcr.io/distroless/static-debian12 AS runner
WORKDIR /app
COPY --from=builder "/app" .
CMD ["/app/service"]
