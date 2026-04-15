# NOTE: build binary
FROM golang:1.26-bookworm AS builder

RUN apt update && apt upgrade -y && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY . .

RUN  go mod download && go mod verify

RUN mkdir -p /bin
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -v -o /bin/api ./cmd/api/main.go

# NOTE: build frontend
FROM node:18-alpine AS frontend

RUN npm install -g pnpm

WORKDIR /app

COPY package.json pnpm-lock.yaml* ./

RUN pnpm install --frozen-lockfile

# Copy the rest of the application
COPY . .

# Build the application
RUN pnpm run build

FROM scratch

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# force rebuild for this part
ARG BUILD_DATE
LABEL rebuild_trigger=$BUILD_DATE
COPY --from=builder /bin/api /api

COPY --from=builder /app/public/index.html /public/index.html
COPY --from=frontend /app/public/*.bundle.js /public/
COPY --from=frontend /app/public/styles.css /public/styles.css

EXPOSE 8080 

CMD [ "/api" ]
