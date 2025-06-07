# Stage 1: Build stage
FROM golang:1.24-alpine AS build
WORKDIR /app
COPY . .
RUN apk add --no-progress --no-cache gcc musl-dev
COPY go.mod go.sum ./
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -tags musl -ldflags '-extldflags "-static"' -o main .

# Stage 2: Final stage
FROM scratch as release
WORKDIR /app
COPY --from=build /app/main .
ARG CONFIG_FILE=dev.env
COPY ${CONFIG_FILE} /app/${CONFIG_FILE}
ENV CONFIG_FILE=${CONFIG_FILE}
ENTRYPOINT ["./main"]
