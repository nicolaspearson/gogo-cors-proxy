# Go Go Cors Proxy

A simple **Go** proxy which adds CORS headers to an incoming request. This allows your application to execute requests on a resource hosted on a different domain.

## Getting Started

1. Download and [install](https://golang.org/doc/install) Go
2. Setup your $GOROOT and $GOPATH in your bashrc / zshrc, for example:
```
    export GOROOT=/usr/local/go
    export GOPATH=/dev/go
    export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
```

## Application Flags

| Flag | Default | Description |
| -------- |:------- |:----------- |
| `target` | localhost:8080 | host:port to proxy requests to |
| `listen` | localhost:8181 | host:port to listen on |
| `protocol` | http | protocol used by the target |
| `host` | localhost:3000 | host header to be used for the proxy request |
| `origin` | `http://localhost:3000` | origin header to be used for the proxy request |
| `methods` | true | enable / disable default access control methods |
| `debug` | false | enable / disable debug messages |

## Running The Application

Once **Go** has been correctly installed and configured, execute:

`go run proxy.go -target=0.0.0.0:8080 -listen=0.0.0.0:8181 -host=localhost:3000 -origin=http://localhost:3000`

Now all incoming requests on port `8181` will be proxied to `http://0.0.0.0:8080`

## Docker

Ensure you have Docker installed before continuing.

### Using The Docker Image

Please see the `docker-compose.yml` file an example of how the image may be used.

### Building The Docker Image

To build the docker image you can simply run: `docker-compose up`

## Contribution Guidelines

### Pull Requests

Here are some basic rules to follow to ensure timely addition of your request:

  1. Match coding style (braces, spacing, etc.).
  2. If it is a feature, bugfix, or anything please only change the minimum amount of code required to satisfy the change.
  3. Please keep PR titles easy to read and descriptive of changes, this will make them easier to merge :)
  4. Pull requests _must_ be made against `develop` branch. Any other branch (unless specified by the maintainers) will get rejected.
  5. Check for existing issues first, before filing a new issue.

## License

MIT License
