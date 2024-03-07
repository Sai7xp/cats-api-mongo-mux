## GoLang Web Server - Cat Manager Service

a web server developed in GoLang Uses MongoDB for data storage and Dockerized for easy deployment

### Usage

Clone the repository:

```bash
git clone https://github.com/your-username/cat-manager.git
```

Build the Docker image:

```bash
docker build -t cats-api .
```

Run the Docker container:

```bash
docker run -p 6000:6000 cats-api
```

Run in detached Mode -d

```bash
# docker run -d -p 6000:6000 cats-api
docker run -p 6000:6000 --link mongodatabase:mongo cats-api
```
Docker run command takes an image and builds the container out of it and runs it

### Dependencies:

- GoLang
- Gorilla Mux
- MongoDB
- Docker
