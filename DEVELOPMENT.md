# flow Development

`flow` is written in [Go](https://golang.org/).

The following `make` commands are available:

|                                | Make Command      |
|--------------------------------|-------------------|
| **Install Local Dependencies** | `make local/deps` |
| **Build**                      | `make go/build`   |
| **Test**                       | `make go/test`    |
| **Pre-commit**                 | `make pre-commit` |

# Installing via Source

```bash
$ git clone github.com/jahvon/flow
$ cd flow
$ go generate ./...
$ go install
```
