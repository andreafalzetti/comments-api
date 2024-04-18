# Contributing

## Getting Started

The easiest way to start contributing is using [Gitpod](https://gitpod.io). The project and its dependencies are configured in the `.gitpod.yml` and `.gitpod.Dockerfile`.

Use the following commands:

- `make build`: build the binary
- `make run`: builds and executes the binary
- `make dev`: builds and watches the code for changes

## Configuration

The service uses the following environment variables as config:

- `COMMENTS_API_HOST` (default: `127.0.0.1`)
- `COMMENTS_API_PORT` (default: `14000`)

## What if I don't have Gitpod?

You can also run this project locally. You need `go`, `ngrok` and `make` installed.

Run:

```bash
make build
make dev
```