# Build frontend
FROM node:20-slim AS frontend-builder
WORKDIR /frontend

# Copy package files first
COPY ./frontend/package*.json ./

# Install dependencies
RUN npm install

# Copy frontend files
COPY ./frontend/ ./

# Run the build
RUN npm run build
RUN echo "=== Frontend Build Output ===" && ls -la build/

# Build backend
FROM golang:1.23.4 AS backend-builder
WORKDIR /build
COPY ./backend/go.mod ./backend/go.sum ./
RUN go mod download
COPY ./backend/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o server

# Final stage - using debian for debugging
FROM debian:bullseye-slim
WORKDIR /app
COPY --from=backend-builder /build/server /app/server
COPY --from=frontend-builder /frontend/build /app/static
RUN echo "=== Static Files in Container ===" && ls -la /app/static

EXPOSE 8080
CMD ["/app/server"]