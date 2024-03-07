# Use the official GoLang image
FROM golang:1.22.0

# Create new directory /app
RUN mkdir /app

# Add all the files to /app folder
ADD . /app

# Set the working directory inside the container
WORKDIR /app

# Download Dependencies
RUN go mod download

# Build
RUN go build -o main

# Expose the port that your GoLang application listens on
EXPOSE 6000

# finally run the executable
CMD [ "/app/main" ]