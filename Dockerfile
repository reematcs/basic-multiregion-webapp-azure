# Build frontend
FROM node:20-slim AS frontend-builder
WORKDIR /frontend

# Copy package files first
COPY ./health-dashboard/frontend/package*.json ./

# Install dependencies
RUN npm install

# Copy frontend files
COPY ./health-dashboard/frontend/ ./

# Run the build
RUN npm run build
RUN echo "=== Frontend Build Output ===" && ls -la build/

# Build backend
FROM golang:1.23.4
WORKDIR /app
COPY ./health-dashboard/backend/go.mod ./health-dashboard/backend/go.sum ./
RUN go mod download && go install github.com/go-delve/delve/cmd/dlv@latest
COPY ./health-dashboard/backend/ ./backend/
COPY ./health-dashboard/frontend/ ./frontend/

EXPOSE 8080 2345
CMD ["dlv", "debug", "--headless", "--listen=:2345", "--api-version=2", "--accept-multiclient", "./backend/main.go"]

# Final stage - using debian for debugging
FROM debian:bullseye-slim
WORKDIR /app
COPY --from=backend-builder /build/server /app/server
COPY --from=frontend-builder /frontend/build /app/static
RUN echo "=== Static Files in Container ===" && ls -la /app/static

EXPOSE 8080
CMD ["/app/server"]