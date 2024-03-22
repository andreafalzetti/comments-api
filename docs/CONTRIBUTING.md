# Contributing

## Getting Started

The easiest way to start contributing is using Gitpod. The project and its dependencies are configured in the `.gitpod.yml` and `.gitpod.Dockerfile`.

Use the following commands:

- `make build`: build the binary
- `make run`: builds and executes the binary
- `make dev`: builds and watches the code for changes

## Configuration

The service uses the following environment variables as config:

- `COMMENTS_API_HOST`
- `COMMENTS_API_PORT`
