package main

import (
	"flag"
	"io"
	"os"
	"runtime"

	phroxy "github.com/dutchsec/ares"

	"github.com/BurntSushi/toml"
	logging "github.com/op/go-logging"
)

var version = "0.1"

var format = logging.MustStringFormatter(
	"%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}",
)

var log = logging.MustGetLogger("ares")

var configFile string

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.StringVar(&configFile, "config", "config.toml", "specifies the location of the config file")
}

func main() {
	flag.Parse()

	var (
		err error
	)

	c := phroxy.New()
	if _, err = toml.DecodeFile(configFile, &c); err != nil {
		panic(err)
	}

	logBackends := []logging.Backend{}
	for _, log := range c.Logging {

		var output io.Writer = os.Stdout
		switch log.Output {
		case "stdout":
		case "stderr":
			output = os.Stderr
		default:
			output, err = os.OpenFile(log.Output, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		}

		if err != nil {
			panic(err)
		}

		backend1 := logging.NewLogBackend(output, "", 0)
		backend1Formatter := logging.NewBackendFormatter(backend1, format)
		backend1Leveled := logging.AddModuleLevel(backend1Formatter)

		level, err := logging.LogLevel(log.Level)
		if err != nil {
			panic(err)
		}

		backend1Leveled.SetLevel(level, "")

		logBackends = append(logBackends, backend1Leveled)
	}

	logging.SetBackend(logBackends...)

	c.ListenAndServe()
}
