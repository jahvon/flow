# flow - Local, CLI Workflow Manager

<p align="center">
  <a href="https://github.com/jahvon/flow/releases" rel="nofollow">
    <img alt="GitHub release (latest SemVer including pre-releases)" src="https://img.shields.io/github/v/release/jahvon/flow?include_prereleases">
  </a>

  <a href="https://github.com/jahvon/flow/actions/workflows/release.yaml" rel="nofollow">
    <img src="https://github.com/jahvon/flow/actions/workflows/release.yaml/badge.svg" alt="goreleaser" style="max-width:100%;">
  </a>

  <a href="https://pkg.go.dev/github.com/jahvon/flow" rel="nofollow">
    <img src="https://pkg.go.dev/badge/github.com/jahvon/flow.svg" alt="Go reference" style="max-width:100%;">
  </a>

  <a href="https://github.com/jahvon/flow/blob/main/LICENSE" rel="nofollow">
    <img src="https://img.shields.io/badge/license-Apache 2.0-blue.svg" alt="License Apache 2.0" style="max-width:100%;">
  </a>

  <br/>

  <a href="https://codecov.io/gh/jahvon/flow" >
    <img src="https://codecov.io/gh/jahvon/flow/branch/main/graph/badge.svg?token=CLP6KW4QLK"/>
  </a>

  <a href="https://github.com/jahvon/flow/actions/workflows/codeql.yaml" rel="nofollow">
    <img src="https://github.com/jahvon/flow/actions/workflows/codeql.yaml/badge.svg" alt="codeql" style="max-width:100%;">
  </a>

  <a href="https://goreportcard.com/report/github.com/jahvon/flow" rel="nofollow">
    <img src="https://goreportcard.com/badge/github.com/jahvon/flow" alt="Go report card" style="max-width:100%;">
  </a>
</p>
<br/>

## Development

### Dependencies

`flow` is written in [Go](https://golang.org/).

Prerequisites:

- [Go 1.20+](https://golang.org/doc/install)

Other local dependencies can be installed with:

```sh
make local/deps
```

### Build / Testing

```sh
make go/build
```

```sh
make go/test
```

### Pre-commit

```sh
make pre-commit
```

## Install

You can install the pre-compiled binary (in several ways), use Docker or compile from source (when on OSS).

Below you can find the steps for each of them.

<details>
  <summary><h3>homebrew tap</h3></summary>

```bash
brew install jahvon/tap/flow
```

</details>

<details>
  <summary><h3>apt</h3></summary>

```bash
echo 'deb [trusted=yes] https://apt.fury.io/jahvon/ /' | sudo tee /etc/apt/sources.list.d/jahvon.list
sudo apt update
sudo apt install flow
```

</details>

<details>
  <summary><h3>yum</h3></summary>

```bash
echo '[jahvon]
name=Gemfury jahvon repository
baseurl=https://yum.fury.io/jahvon/
enabled=1
gpgcheck=0' | sudo tee /etc/yum.repos.d/jahvon.repo
sudo yum install goreleaser
```

</details>

<details>
  <summary><h3>deb, rpm and apk packages</h3></summary>
Download the .deb, .rpm or .apk packages from the [release page](https://github.com/jahvon/flow/releases) and install them with the appropriate tools.
</details>

<details>
  <summary><h3>go install</h3></summary>

```bash
go install github.com/jahvon/flow@latest
```

</details>

<details>
  <summary><h3>from the GitHub releases</h3></summary>

Download the pre-compiled binaries from the [release page](https://github.com/jahvon/flow/releases) page and copy them to the desired location.

```bash
$ VERSION=v1.0.0
$ OS=Linux
$ ARCH=x86_64
$ TAR_FILE=flow_${OS}_${ARCH}.tar.gz
$ wget https://github.com/jahvon/flow/releases/download/${VERSION}/${TAR_FILE}
$ sudo tar xvf ${TAR_FILE} flow -C /usr/local/bin
$ rm -f ${TAR_FILE}
```

</details>

<details>
  <summary><h3>manually</h3></summary>

```bash
$ git clone github.com/jahvon/flow
$ cd flow
$ go generate ./...
$ go install
```

</details>

## Credits

Repository initially generated with [thazelart/golang-cli-template](https://github.com/thazelart/golang-cli-template).
