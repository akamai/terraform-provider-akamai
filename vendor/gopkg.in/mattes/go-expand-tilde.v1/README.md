# Expand ~ in path

[![Build Status](https://travis-ci.org/mattes/go-expand-tilde.svg?branch=v1)](https://travis-ci.org/mattes/go-expand-tilde)
[![GoDoc](https://godoc.org/gopkg.in/mattes/go-expand-tilde.v1?status.svg)](https://godoc.org/gopkg.in/mattes/go-expand-tilde.v1)

```
go get gopkg.in/mattes/go-expand-tilde.v1
```

```go
import (
  "gopkg.in/mattes/go-expand-tilde.v1"
)

func main() {
  path, err := tilde.Expand("~/path/to/whatever")
}
```