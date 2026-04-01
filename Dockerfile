FROM golang:1.22-alpine

WORKDIR /app

# Install git for fetching Go modules from VCS if needed.
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

EXPOSE 8555

CMD ["go", "run", "main.go"]
