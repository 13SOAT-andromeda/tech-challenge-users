FROM golang:1.23-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/api

FROM gcr.io/distroless/static-debian12 AS production

COPY --from=build /app/server /server

EXPOSE 8081

ENTRYPOINT ["/server"]
