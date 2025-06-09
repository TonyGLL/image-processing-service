# Image Processing Service

This is a Go-based image processing service that provides a RESTful API for managing users and (eventually) processing images. It uses JWT for authentication and PostgreSQL as its database.

## Features

*   **User Management:** Create and manage user accounts.
*   **Authentication:** Secure API endpoints using JSON Web Tokens (JWT).
*   **Database:** Uses PostgreSQL for data storage.
*   **API Documentation:** Provides Swagger documentation for the API.
*   **Containerized:** Uses Docker for development and deployment.
*   **Live Reloading:** Uses `air` for live reloading during development.

## Roadmap

Idea based on an intermediate Golang project from the roadmap.sh site:
[https://roadmap.sh/projects/image-processing-service](https://roadmap.sh/projects/image-processing-service)

## Prerequisites

Before you begin, ensure you have the following installed:

*   [Go](https://golang.org/doc/install) (version 1.24.4 or newer)
*   [Docker](https://docs.docker.com/get-docker/)
*   [Docker Compose](https://docs.docker.com/compose/install/)
*   [migrate](https://github.com/golang-migrate/migrate#installation) (specifically, `golang-migrate/migrate`)

## Getting Started

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/TonyGLL/image-processing-service.git
    cd image-processing-service
    ```

2.  **Configuration:**
    The application requires a configuration file. Create a file (e.g., `local.env`) in the root directory with the following content. **Remember to replace placeholder values, especially `JWT_SECRET_KEY`.**
    ```env
    DB_DRIVER=postgres
    DB_SOURCE=postgresql://root:secret@localhost:5432/image_processing?sslmode=disable
    SERVER_ADDRESS=0.0.0.0:8080
    JWT_SECRET_KEY=your-very-strong-and-secret-jwt-key # Replace with a strong secret key
    ACCESS_TOKEN_DURATION=15m
    REFRESH_TOKEN_DURATION=24h
    ```
    After creating the file, set the `CONFIG_FILE` environment variable to its name:
    ```bash
    export CONFIG_FILE=local.env
    ```
    (You might want to add this export to your shell's configuration file, like `.bashrc` or `.zshrc`, for persistence).

3.  **Start PostgreSQL Container:**
    This command starts a PostgreSQL database instance in a Docker container.
    ```bash
    make postgres
    ```

4.  **Create the Database:**
    This command connects to the running PostgreSQL container and creates the `image_processing` database.
    ```bash
    make createdb
    ```

5.  **Run Database Migrations:**
    This command applies the necessary SQL schema migrations to the database.
    ```bash
    make migrateup
    ```

6.  **Build and Start the Application:**
    This command builds the Go application and then starts it using `air` for live reloading.
    ```bash
    make start
    ```
    The application should now be accessible at `http://localhost:8080` (or the `SERVER_ADDRESS` defined in your `local.env` file).

## API Endpoints

The API is versioned under `/api/v1`.

*   **Authentication:**
    *   `POST /auth/login`: Login a user and receive JWT tokens. (Body: `{"username": "youruser", "password": "yourpassword"}`)
*   **Users:**
    *   `POST /users`: Create a new user. (Requires JWT authentication in the `Authorization` header: `Bearer <token>`) (Body: `{"username": "newuser", "password": "newpassword", "email": "user@example.com"}`)

For detailed and interactive API documentation, navigate to `/swagger/index.html` in your browser once the service is running.

## Makefile Commands

The `Makefile` provides several commands to streamline development:

*   `make postgres`: Starts the PostgreSQL Docker container.
*   `make stop-postgres`: Stops the PostgreSQL Docker container.
*   `make createdb`: Creates the `image_processing` database within the PostgreSQL container.
*   `make dropdb`: Drops the `image_processing` database.
*   `make migrateup`: Applies all pending database migrations to the configured database.
*   `make migratedown`: Reverts the last applied database migration.
*   `make build`: Compiles the Go application, placing the output binary in the root directory.
*   `make start`: Builds the application and then runs it using `air` for live-reloading (useful for development).
*   `make initschema NAME=your_migration_name`: Creates a new SQL migration file in `sql/migrations/`. Replace `your_migration_name` with a descriptive name (e.g., `add_image_table`).

## Configuration Details

Application configuration is loaded from an environment file (e.g., `local.env`) specified by the `CONFIG_FILE` environment variable.

Key configuration variables:

*   `DB_DRIVER`: Database driver (currently supports `postgres`).
*   `DB_SOURCE`: Database connection string (e.g., `postgresql://user:password@host:port/dbname?sslmode=disable`).
*   `SERVER_ADDRESS`: Host and port for the API server (e.g., `0.0.0.0:8080`).
*   `JWT_SECRET_KEY`: **Critical:** A strong, unique secret key for signing JWT tokens.
*   `ACCESS_TOKEN_DURATION`: Validity duration for access tokens (e.g., `15m`, `1h`).
*   `REFRESH_TOKEN_DURATION`: Validity duration for refresh tokens (e.g., `24h`, `7d`).

## Project Structure

```
.
├── Dockerfile            # Defines the Docker image for the application.
├── Makefile              # Contains build, run, and utility commands.
├── README.md             # This file.
├── air.toml              # Configuration for 'air' (live-reloading tool for Go).
├── docker-compose.yml    # Defines services, networks, and volumes for Docker.
├── go.mod                # Go module definition file.
├── go.sum                # Go module checksums.
├── internal/             # Directory for internal application code, not intended for import by other projects.
│   ├── api/              # Handles API requests, routing, and server logic.
│   ├── db/               # Database interaction layer, including SQLC-generated code.
│   ├── middlewares/      # HTTP middleware (e.g., JWT authentication).
│   └── utils/            # Utility functions (config loading, JWT generation, etc.).
├── main.go               # Entry point of the application.
├── sql/                  # Contains SQL files.
│   ├── migrations/       # Database schema migration files (managed by golang-migrate).
│   └── queries/          # SQL queries used by SQLC to generate Go database interface code.
└── sqlc.yml              # Configuration file for SQLC (SQL compiler).
```

## Contributing

Contributions are welcome! If you have suggestions or improvements, please open an issue or submit a pull request.

## License

The license for this project is yet to be determined.
