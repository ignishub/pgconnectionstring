# pgconnectionstring

![Go](https://github.com/ignishub/pgconnectionstring/workflows/Go/badge.svg) [![Go Report Card](https://goreportcard.com/badge/github.com/ignishub/pgconnectionstring)](https://goreportcard.com/report/github.com/ignishub/pgconnectionstring)

This library provides utility function that parses connection string for github.com/lib/pq.
Ths is modified version of code from this library.

```go

values, err := pgconnectionstring.Parse("user=username password='pwd string'")
```
This functions return values as string map. Returned values can be converted to connection string:

```go
ToConnectionString(values)
```