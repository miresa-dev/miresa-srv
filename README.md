<h1 align="center">Miresa</h1>

<p align="center">
Miresa is minimal, free software for building online forums and communities.
<br/><br/>
<a href="./LICENSE">
    <img
        alt="License: Unlicense"
        src="https://img.shields.io/badge/License-Unlicense-green.svg"
    />
</a>
<a href="https://goreportcard.com/github.com/miresa-dev/miresa-srv">
    <img
        alt="Go Report Card"
        src="https://goreportcard.com/badge/github.com/miresa-dev/miresa/srv"
    />
</a>
<a href="https://github.com/miresa-dev/miresa-srv/issues">
    <img
        alt="Open issues"
        src="https://img.shields.io/github/issues/miresa-dev/miresa-srv"
    />
</a>
<a href="https://github.com/miresa-dev/miresa-srv/actions/workflows/golangci-lint.yml">
    <img
        alt="Linter status"
        src="https://github.com/miresa-dev/miresa-srv/actions/workflows/golangci-lint.yml/badge.svg"
    />
</a>
</p>

**NOTE**: Miresa Server is still in the pre-alpha development stage, suitable
for testing purposes only.

## Demo

Insert images here

## Installation

If you want to customize the CSS, Miresa nees to be run from the root of the
repository. For this, you should clone the repo:

```bash
git clone https://github.com/miresa-dev/miresa-srv
````

Then `cd` into the directory and build the project:

```bash
cd miresa-srv
go build .
```

If you're fine with the defaults, you can install with Go:

```bash
go install github.com/miresa-dev/miresa-srv@latest
``` 

Or get a binary form the
[releases page](https://github.com/miresa-dev/miresa-srv/releases).

## Usage

This section is for administrators. If you're a user, you may want to check out
the [official CLI](https://github.com/miresa-dev/mirec) or the
[fancier web client](https://github.com/miresa-dev/mirer). The built-in web
client is extremely minimal.

You can start the server by running `miresa-srv`:

```bash
./miresa-srv
# Or, if you installed with `go install`,
miresa-srv
```

See the [selfhosting docs](https://miresa-dev.github.io/doc/selfhost) for more
in-depth information.

## Support

If you need support, you can open a GitHub discussion, send a message on
[Gitter](https://matrix.to/#/#miresa:gitter.im), or send us an email.

## Roadmap

* Support the [full API](https://miresa-dev.github.io/doc/api/ref)
* Make the web client
* Users
  * [x] Create
    * [ ] Validate session ID and captcha
  * [x] Read
  * [ ] Update
  * [ ] Delete
* Items
  * [ ] Create
  * [ ] Read
  * [ ] Update
  * [ ] Delete
* Configuration
  * [ ] Allow JSON/YAML configuration
  * [ ] More config options
    * [ ] ID length
    * [ ] What info to show on `/v`
      * [x] Goroutine count
      * [x] OS
      * [x] Arch
      * [x] Current server time
      * [ ] Uptime
* Rate-limiting
  * [ ] 30 requests per minute

## Contributing

All sorts of contributions are always welcome! See the [contribution docs](https://miresa-dev.github.io/doc/code/contrib) for ways to help.

## Acknowledgements

### Contributors

* [Kaamkiya](https://github.com/Kaamkiya)
<!--S:CONTRIBUTORS-->
<!--E:CONTRIBUTORS-->

### Libraries

* [Chi](https://go-chi.io)

## License

This project is licensed under the [Unlicense](./LICENSE). All code is public domain unless otherwise specified.

