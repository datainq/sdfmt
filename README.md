# sdfmt [![Build Status](https://travis-ci.org/datainq/sdfmt.svg?branch=master)](https://travis-ci.org/datainq/sdfmt)

Stackdriver formatter for [logrus](https://github.com/sirupsen/logrus) logger compliant with [Stackdriver Logging API](https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry).

## Installation

`go get -u github.com/datainq/sdfmt`

## Usage

```go
package main

import (
	"github.com/datainq/sdfmt"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&sdfmt.StackdriverFormatter{})
	log.WithFields(log.Fields{
		"animal": "walrus",
	}).Info("A walrus appears")
}
```

Output:

`{"labels":{"animal":"walrus"},"message":"A walrus appears","severity":200,"timestamp":"2019-03-31T18:35:13.752104+02:00"}`

## Logging levels association

|              `logrus.Level`              | Stackdriver's `LogSeverity` |
|:----------------------------------------:|:---------------------------:|
|                     -                    |         DEFAULT (0)         |
| `logrus.TraceLevel`, `logrus.DebugLevel` |         DEBUG (100)         |
|            `logrus.InfoLevel`            |          INFO (200)         |
|                     -                    |         NOTICE (300)        |
|            `logrus.WarnLevel`            |        WARNING (400)        |
|            `logrus.ErrorLevel`           |         ERROR (500)         |
|                     -                    |        CRITICAL (600)       |
|                     -                    |         ALERT (700)         |
| `logrus.FatalLevel`, `logrus.PanicLevel` |       EMERGENCY (800)       |
