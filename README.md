# Website Health Checker

Golang website and server designed to Health-Check a user-defined list of URLs/websites and display them in a status Dashboard.

Technologies within this repository:

- Go (1.26.2): Programming Language

- Templ (0.3.1020): Tool to embed Go in HTML templates

- HTMX (2.0.10): Library to build modern interactive interfaces directly in HTML

- Tailwind CSS (V4.3): Utility-first CSS framework

- Air (github.com/air-verse/air v1.65.1): Dynamic reloading tool

- Nix: Reproducible and declarative package management tool

- Docker: ...

- Docker Compose: ...

- Bruno: ...

- Prometheus: ...

- Grafana: ...

> Disclaimer: This project used AI for guiding, learning, and controlled code generation. No decision were made without human analysis and intervention.

## Installation and Usage

To proceed, you will need **Go** installed, alongside its utils. You can install it via [Nix](https://nixos.org/download/) (recommended), your package manager (apt, pacman, winget, etc.), or downloading the binary on their [website](https://go.dev/doc/install).

###### Installing Go with Nix

``` zsh
nix develop # Sandboxed; `exit` to escape
```

#### Running the server

Default (recommended)
``` zsh
go run .
```

For development (dynamic reloading; devs-only)
``` zsh
go tool air
```

#### Using the app

Open http://localhost:8080 on your browser.

TODO...


## Tests

To properly test the software, the following tools are used:

`go vet` - Identify possible issues and bugs.

`go test -race` - Run several tests based on the [Test Plan](docs/test_plan.md), including race conditions (`race` flag).

#### Test Plan

TODO...

Check the [docs](docs/test_plan.md) for more details.

## Folder Structure

TODO...

## Internal Structure

TODO...

## License

This project in under an [MIT License](LICENSE).
