# assert [![GoDoc](https://godoc.org/github.com/tj/assert?status.svg)](https://godoc.org/github.com/tj/assert)

Assertion pkg for Go, copied from [github.com/stretchr/testify's](//github.com/stretchr/testify) require package.

I find early errors more useful than the `t.Errorf()` calls, which often fall through to nil pointers causing panics etc.
