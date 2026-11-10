FROM golang:1.25-alpine AS build
WORKDIR /src
COPY . .
# go.mod in the repo predates the pgx dependency, so resolve it at build time.
RUN go mod tidy && go mod download
# The committed generated.go does not compile; regenerate it from the schema.
RUN go run github.com/99designs/gqlgen generate
RUN CGO_ENABLED=0 go build -o /out/server .

FROM alpine:3.20
RUN adduser -D app
WORKDIR /app
COPY --from=build /out/server ./server
COPY migrations ./migrations
USER app
# PORT is provided by the host (Render, Fly, etc.); the server falls back to 8080.
EXPOSE 8080
CMD ["./server"]
