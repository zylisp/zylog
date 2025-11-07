# zylog

*A simple wrapper for customized logrus usage*


## Setup

Here's an example of setting up global logging for use by your app in your
app's `logging` package, based on configuration pulled in by Viper (from either
a config file or ENV variables):


```go
package logging

import (
	cfg "github.com/spf13/viper"
	log "github.com/geomyidia/zylog/logger"
)

func init() {
	log.SetupLogging(&log.ZyLogOptions{
		Colored:         cfg.GetBool("logging.colored"),
		Level:           cfg.GetString("logging.level"),
		Output:          cfg.GetString("logging.output"),
		ReportCaller:    cfg.GetBool("logging.report-caller"),
		TimestampFormat: log.Simple, // Optional: RFC3339, Simple, or Time
	})
}
```


## Usage

The setup configures the logrus logger, so wherever you want to log, simply
use logrus as you would normally:

```go
...
import (
	log "github.com/sirupsen/logrus"
)
...
log.Info("You are standing in an open field west of a white house.")
...
```


There's some more example usage in the demo (`./cmd/zylog-demo/main.go`). To run it:

```bash
$ make build
$ ./bin/demo
```

At which point you should see something like the following:

![screenshot](assets/images/screenshot.png)


## Timestamp Formats

Zylog supports three timestamp formats:

- `logger.RFC3339` - Full RFC3339 format (e.g., `2025-11-07T14:30:45-08:00`)
- `logger.Simple` - Compact format: YYYYMMDD.HHmmSS (e.g., `20251107.143045`) **[Default]**
- `logger.Time` - Time only: HH:mm:SS (e.g., `14:30:45`)

### Examples

```go
// Using Simple format (default)
log.SetupLogging(&log.ZyLogOptions{
	Colored:         true,
	Level:           "info",
	Output:          "stdout",
	ReportCaller:    false,
	TimestampFormat: log.Simple,
})

// Using RFC3339 format
log.SetupLogging(&log.ZyLogOptions{
	Colored:         true,
	Level:           "info",
	Output:          "stdout",
	ReportCaller:    false,
	TimestampFormat: log.RFC3339,
})

// Using Time-only format
log.SetupLogging(&log.ZyLogOptions{
	Colored:         true,
	Level:           "info",
	Output:          "stdout",
	ReportCaller:    false,
	TimestampFormat: log.Time,
})

// If TimestampFormat is not specified, Simple is used by default
log.SetupLogging(&log.ZyLogOptions{
	Colored:      true,
	Level:        "info",
	Output:       "stdout",
	ReportCaller: false,
})
```


## Development

A convenience Bash environment file is provided for easy setup:

```bash
$ . .env
```


## Background

Note that the formatting provided by this util lib is inspired by the
[Twig Clojure](https://github.com/clojusc/twig) and the
[Logjam LFE](https://github.com/lfex/logjam) libraries.


## License

© 2019, ZYLISP Project

© 2019, geomyidia Project

Apache License, Version 2.0
