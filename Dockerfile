FROM golang:1.25-alpine

WORKDIR /app


COPY go.mod go.sum ./
RUN go mod download

EXPOSE 8080

COPY . .
ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target="/root/.cache/go-build" go build -o srv cmd/api/main.go


CMD ["/app/srv"]
