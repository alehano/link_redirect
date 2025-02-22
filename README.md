# Link Redirection Service

This is a simple Go-based service that redirects incoming links based on a configuration file. It uses the `chi` router for handling HTTP requests and `configor` for dynamic configuration management.

## Features

- Redirects incoming links to specified URLs.
- Configuration is stored in a `config.yml` file.
- Automatically reloads configuration every 10 seconds.

## Requirements

- Go 1.16 or later
- `github.com/go-chi/chi/v5` for server routing
- `github.com/jinzhu/configor` for dynamic configuration

## Setup

1. **Clone the repository:**

   ```bash
   git clone <repository-url>
   cd <repository-directory>
   ```

2. **Install dependencies:**

   ```bash
   go get github.com/go-chi/chi/v5
   go get github.com/jinzhu/configor
   ```

3. **Create a `config.yml` file:**

   ```yaml
   urls:
     - url: link1
       redirect: https://example.com/123
     - url: link2
       redirect: https://example.com/456
   ```

4. **Run the application:**

   ```bash
   go run main.go
   ```

## Usage

- Start the server and access it via `http://localhost:8080/{link}`.
- The server will redirect to the URL specified in the `config.yml` file for the given `{link}`.

## Configuration

- The `config.yml` file contains the URLs and their corresponding redirection targets.
- The configuration is reloaded every 10 seconds, allowing for dynamic updates without restarting the server.

## Notes

- Ensure that the `config.yml` file is in the same directory as `main.go`.
- The service listens on port 8080 by default.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
