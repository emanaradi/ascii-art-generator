# Defines how to build the application into a Docker image
# This happens in multiple steps: 
    # First, get GO
    # Bild the go application
    # Create a smaller container
    # Copy the application and its required files
    # Run the go server

# 1. Define the builder image
FROM golang:1.23-alpine AS builder

#  2. Set the working directory inside the container
WORKDIR /app

# 3. Copy Go module files
COPY go.mod ./

# 4. Download dependencies
RUN go mod download

# 5. Copy the rest of the project
COPY . .

# 6. Compile the Go application
RUN go build -o ascii-art-web-dockerize ./main

# 7. Define the final runtime image
FROM alpine:latest

# 8.Add metadata to the Docker image
LABEL org.opencontainers.image.title="ASCII Art Web Dockerize"
LABEL org.opencontainers.image.description="A web application for generating ASCII art."
LABEL org.opencontainers.image.version="1.0.0"
LABEL org.opencontainers.image.authors="Eman Alaradi, Noor Alhayki"

# 9.Set the working directory
WORKDIR /app

# 10.Copy the compiled application
COPY --from=builder /app/ascii-art-web-dockerize .

# 11. Copy files required by the application
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/banners ./banners
COPY --from=builder /app/assets ./assets

# 12. Set the port used by the application
EXPOSE 8081

# 13. Finally, start the application when the container runs
CMD ["./ascii-art-web-dockerize"]