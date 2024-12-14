# Build frontend
FROM node:18-slim AS frontend-builder
WORKDIR /frontend

# Install CLI tools needed for shadcn-ui and craco
RUN npm install -g @shadcn/ui @craco/craco

# Copy package files first
COPY ./frontend/package*.json ./

# Install dependencies including craco
RUN npm install --legacy-peer-deps \
    @craco/craco \
    @radix-ui/react-alert-dialog \
    @radix-ui/react-slot \
    class-variance-authority \
    clsx \
    tailwindcss \
    tailwindcss-animate \
    tailwind-merge \
    lucide-react \
    recharts

# Create necessary directories
RUN mkdir -p src/components/ui src/lib public

# Copy frontend files
COPY ./frontend/ ./

# Configure craco for path aliases - now in proper JS module format
RUN echo 'module.exports = {\
  webpack: {\
    alias: {\
      "@": require("path").resolve(__dirname, "src")\
    }\
  }\
};' > craco.config.js

# Update package.json scripts to use craco
RUN npm pkg set scripts.build="craco build"

# Add utils file for cn function
RUN echo 'import { clsx } from "clsx";\
import { twMerge } from "tailwind-merge";\
\
export function cn(...inputs) {\
  return twMerge(clsx(inputs));\
}' > src/lib/utils.js

# Add shadcn component files
RUN echo 'import * as React from "react";\
const Card = React.forwardRef(({ className, ...props }, ref) => (\
  <div\
    ref={ref}\
    className="rounded-lg border bg-card text-card-foreground shadow-sm"\
    {...props}\
  />\
));\
\
const CardHeader = React.forwardRef(({ className, ...props }, ref) => (\
  <div\
    ref={ref}\
    className="flex flex-col space-y-1.5 p-6"\
    {...props}\
  />\
));\
\
const CardTitle = React.forwardRef(({ className, ...props }, ref) => (\
  <h3\
    ref={ref}\
    className="text-2xl font-semibold leading-none tracking-tight"\
    {...props}\
  />\
));\
\
const CardContent = React.forwardRef(({ className, ...props }, ref) => (\
  <div ref={ref} className="p-6 pt-0" {...props} />\
));\
\
export { Card, CardHeader, CardTitle, CardContent };' > src/components/ui/card.js

RUN echo 'import * as React from "react";\
const Alert = React.forwardRef(({ className, variant = "default", ...props }, ref) => (\
  <div\
    ref={ref}\
    role="alert"\
    className="relative w-full rounded-lg border p-4"\
    {...props}\
  />\
));\
\
const AlertTitle = React.forwardRef(({ className, ...props }, ref) => (\
  <h5\
    ref={ref}\
    className="mb-1 font-medium leading-none tracking-tight"\
    {...props}\
  />\
));\
\
const AlertDescription = React.forwardRef(({ className, ...props }, ref) => (\
  <div\
    ref={ref}\
    className="text-sm"\
    {...props}\
  />\
));\
\
export { Alert, AlertTitle, AlertDescription };' > src/components/ui/alert.js

RUN echo 'import * as React from "react";\
const Button = React.forwardRef(({ className, variant = "default", size = "default", ...props }, ref) => (\
  <button\
    ref={ref}\
    className="inline-flex items-center justify-center rounded-md text-sm font-medium transition-colors focus:outline-none focus:ring-2 focus:ring-offset-2 disabled:opacity-50 disabled:pointer-events-none"\
    {...props}\
  />\
));\
\
export { Button };' > src/components/ui/button.js

# Add index.html
RUN echo '<!DOCTYPE html>\
<html lang="en">\
  <head>\
    <meta charset="utf-8" />\
    <meta name="viewport" content="width=device-width, initial-scale=1" />\
    <title>Health Dashboard</title>\
  </head>\
  <body>\
    <div id="root"></div>\
  </body>\
</html>' > public/index.html

# Add Tailwind config
RUN echo 'module.exports = {\
  content: ["./src/**/*.{js,jsx}"],\
  theme: {\
    extend: {},\
  },\
  plugins: [require("tailwindcss-animate")],\
}' > tailwind.config.js

# Run the build
RUN npm run build

# Build backend
FROM golang:1.23.4 AS backend-builder
WORKDIR /build
COPY ./backend/go.mod ./backend/go.sum ./
RUN go mod download
COPY ./backend/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o server

# Final stage
FROM gcr.io/distroless/static:nonroot
COPY --from=backend-builder /build/server /server
COPY --from=frontend-builder /frontend/build /static
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/server"]


# # Build frontend
# FROM node:18-slim AS frontend-builder
# WORKDIR /frontend

# # Copy package files first
# COPY ./frontend/package*.json ./

# # Install dependencies
# RUN npm install

# # Copy remaining frontend files
# COPY ./frontend/ ./

# # Fix path aliases by adding jsconfig.json if it doesn't exist
# RUN echo '{ \
#   "compilerOptions": { \
#     "baseUrl": ".", \
#     "paths": { \
#       "@/*": ["./src/*"] \
#     } \
#   } \
# }' > jsconfig.json

# # Configure vite.config.js for path aliases
# RUN echo 'import { defineConfig } from "vite"; \
# import react from "@vitejs/plugin-react"; \
# import path from "path"; \
# \
# export default defineConfig({ \
#   plugins: [react()], \
#   resolve: { \
#     alias: { \
#       "@": path.resolve(__dirname, "./src") \
#     } \
#   } \
# });' > vite.config.js

# RUN npm run build

# # Build backend
# FROM golang:1.23.4 AS backend-builder
# WORKDIR /build
# COPY ./backend/go.mod ./backend/go.sum ./
# RUN go mod download
# COPY ./backend/ ./
# RUN CGO_ENABLED=0 GOOS=linux go build -o server

# # Final stage
# FROM gcr.io/distroless/static:nonroot
# WORKDIR /
# COPY --from=backend-builder /build/server .
# COPY --from=frontend-builder /frontend/build ./static
# EXPOSE 8080
# USER nonroot:nonroot
# ENTRYPOINT ["/server"]