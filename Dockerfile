#1-stage binary
FROM golang:latest AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /large-assignment

#2-stage app image
FROM golang:latest AS gateway-image
WORKDIR /app
COPY --from=build /large-assignment /large-assignment
COPY config/config.yaml config/config.yaml
EXPOSE 3000
CMD ["/large-assignment"]