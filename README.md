# hourstable data type

[![Build Status](https://travis-ci.org/geniusrabbit/hourstable.svg?branch=master)](https://travis-ci.org/geniusrabbit/hourstable)
[![Go Report Card](https://goreportcard.com/badge/github.com/geniusrabbit/hourstable)](https://goreportcard.com/report/github.com/geniusrabbit/hourstable)
[![GoDoc](https://godoc.org/github.com/geniusrabbit/hourstable?status.svg)](https://godoc.org/github.com/geniusrabbit/hourstable)
[![Coverage Status](https://coveralls.io/repos/github/geniusrabbit/hourstable/badge.svg)](https://coveralls.io/github/geniusrabbit/hourstable)

Implementation of the hour-table for a week (represents active hours by days of the week).
It could be used as the SQL compatible data-type ad stored as String or JSON type in relational or document-oriented DB like PostgreSQL, MySQL, Oracle, MongoDB, etc.

## Hourls-table example

| Day of week / Hour | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | 10 | 11 | 12 | 13 | 14 | 15 | 16 | 17 | 18 | 19 | 20 | 21 | 22 | 23 | 24 |
|:-------------------|:-:|:-:|:-:|:-:|:-:|:-:|:-:|:-:|:-:|:--:|:--:|:--:|:--:|:--:|:--:|:--:|:--:|:--:|:--:|:--:|:--:|:--:|:--:|:--:|
| Monday             | X | X | X |   |   |   |   |   |   |    |    |    |    |    |    |    |    |    |    |    |    |    | X  | X  |
| Tuesday            | X | X | X | X | X | X | X | X | X |    |    |    |    |    |    |    |    |    |    |    |    |    | X  | X  |
| Wednesday          | X | X | X | X | X | X | X | X | X |    |    |    |    |    |    |    |    |    |    |    |    |    | X  | X  |
| Thursday           | X | X | X | X | X | X | X | X | X |    |    |    |    |    |    |    |    |    |    |    |    |    | X  | X  |
| Friday             | X | X | X | X | X | X | X | X | X |    |    |    |    |    |    |    |    |    |    |    |    |    | X  | X  |
| Saturday           | X | X | X |   |   |   |   |   |   |    |    |    |    |    |    |    |    |    |    |    |    |    | X  | X  |
| Sunday             | X | X | X |   |   |   |   |   |   |    |    |    |    |    |    |    |    |    |    |    |    |    | X  | X  |

In a database, it could be stored as TEXT of '1' and '0' symbols

```
111111111111111111111111111111111111111111111111000000000000000000000000000000000000000000000000000000000000000000000000111111111111111111111111111111111111111111111111
```

or as the JSON object

```json
{
  "mon": "111111111111111111111111",
  "tue": "111111111111111111111111",
  "wed": "000000000000000000000000",
  "thu": "000000000000000000000000",
  "fri": "000000000000000000000000",
  "sat": "111111111111111111111111",
  "sun": "000000001111111111111111"
}
```

or short form of JSON

```json
{
  "mon": "*",
  "tue": "*",
  "sat": "*",
  "sun": "000000001111111111111111"
}
```

## Testing & benchmarks

```sh
go test -timeout 30s github.com/demdxx/hourstable -v -race
```

### Benchmarks

```sh
go test -benchmem -run=^$ github.com/demdxx/hourstable -bench . -v -race

pkg: github.com/demdxx/hourstable
Benchmark_Hours-8   	 1000000	      1776 ns/op	       0 B/op	       0 allocs/op
PASS
coverage: 8.3% of statements
ok  	github.com/demdxx/hourstable	2.826s
Success: Benchmarks passed.
```
