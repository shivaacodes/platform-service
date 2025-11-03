# Build stage
FROM golang:1.25-alpine AS build
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o service ./cmd/service

# Final image
FROM scratch
COPY --from=build /app/service /service
EXPOSE 8080
ENTRYPOINT ["/service"]

