# Book Collection API

A simple Go application for playing around with various technologies (listed below).

1. **Search** for books via an external API (e.g., [Open Library](https://openlibrary.org/developers/api)).
2. **Add** searched books to your personal collection stored in a Postgres database.
3. Perform basic **CRUD** (create, read, update, delete) operations on your local book collection.

## Features

- **Search External API**  
  Make a POST request with a title or ISBN to retrieve book data from an external service.  
- **Create Books**  
  Add retrieved (or manually provided) book data to your local Postgres database.  
- **Retrieve Books**  
  List all books or get details of a specific book by ID.  
- **Update Books**  
  Update details (title, author, publication date, etc.) in your local database.  
- **Delete Books**  
  Remove a book from your collection entirely.

## Tech Stack

- **Language**: [Go (Golang)](https://go.dev/)
- **Database**: [PostgreSQL](https://www.postgresql.org/)
- **Containerization**: [Docker](https://www.docker.com/)
- **Orchestration** (Optional): [Kubernetes](https://kubernetes.io/)
- **CI/CD**: [GitHub Actions](https://docs.github.com/en/actions)
- **Infrastructure/Configuration** (Optional): [Ansible](https://www.ansible.com/) or [Terraform](https://www.terraform.io/)

## Getting Started

### Prerequisites

- **Docker** and **Docker Compose**  
  - Install via [Docker’s documentation](https://docs.docker.com/get-docker/).
- **Go** (optional if you’re only running in Docker)  
  - Install from the [Go website](https://go.dev/dl/), if you plan to build/test locally without Docker.
- **PostgreSQL** (optional if not using Docker for the DB)  
  - If you want to run a local PostgreSQL server outside of Docker.

### Quick Start (Docker Compose)

1. **Clone the repository**:
   ```
   git clone https://github.com/mckaywilkerson/book-collection.git
   cd book-collection 
   ```

2. **Start the services**:
    ``` 
    cd deploy/docker
    docker-compose up --build 
    ```

3. **Access the API**:
    * The Go application should be running on port 8081.
    * try curl http://localhost:8081/health (or any other implemented endpoint).

### Usage

### Testing

### Kubernetes (Optional)

### CI/CD with GitHub Actions

### Infrastructure as Code (Optional)

### Roadmap

### Contributing
Pull requests are welcome, as I am new to all of these technologies. For major changes, please open an issue first to discuss what you would like to change.