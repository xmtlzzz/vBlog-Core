FROM node:20-alpine AS frontend
WORKDIR /app/web
COPY web/package*.json ./
RUN npm ci
COPY web/ ./
RUN npm run build

FROM golang:1.22-alpine AS backend
WORKDIR /app/server
COPY server/go.* ./
RUN go mod download
COPY server/ ./
COPY --from=frontend /app/web/dist ./static
RUN go build -o /vblog ./cmd/main.go

FROM alpine:3.19
RUN apk add --no-cache ca-certificates
COPY --from=backend /vblog /usr/local/bin/vblog
COPY deploy/init.sql /init.sql
EXPOSE 8080
CMD ["vblog"]
