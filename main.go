package main

import (
    "time"
	"fmt"
	"html/template"
    "os"
    "syscall"
    "os/signal"
    "github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/go-kit/kit/log"
	"github.com/ghostsquad/go-timejumper"
    "github.com/oklog/run"
    "github.com/fsnotify/fsnotify"
)

const renderedCards = "index.html"

type CardData struct {
}

func main() {
    clock := timejumper.RealClock{}

    defaultTimestampUTC := log.TimestampFormat(
		func() time.Time { return clock.Now().UTC() },
		time.RFC3339Nano,
	)

    var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", defaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
    }

    
    watcher, err := fsnotify.NewWatcher()
    defer watcher.Close()
    if err != nil {
        logger.Log("event", "exit", "reason", err)
        os.Exit(1)
    }
    
    runGroup := run.Group{}
	{
		cancelInterrupt := make(chan struct{})
		runGroup.Add(
            createSignalWatcher(cancelInterrupt), 
            func(error) {
			    close(cancelInterrupt)
            })
    }
    {
        runGroup.Add(
            func() error {
                
            },
            func(error) {
			   
            })
    }

    t := template.Must(template.New("email.tmpl").Parse(indexTmpl))

	err := t.Execute(os.Stdout, getCardData())
	if err != nil {
        logger.Log("event", "exit", "reason", errors.Wrap(err, "Failed to render template!")
		os.Exit(1)
	}

	logger.Log("event", "exit", "reason", runGroup.Run())
}

func getCardData() *CardData {
	return &CardData{}
}

// This function just sits and waits for ctrl-C
func createSignalWatcher(cancelInterruptChan <-chan struct{}) func() error {
	return func() error {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		select {
		case sig := <-c:
			return fmt.Errorf("received signal %s", sig)
		case <-cancelInterruptChan:
			return nil
		}
	}
}

var indexTmpl = `<!doctype html>
<html>
    <head>
        <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
        <meta http-equiv="cache-control" content="max-age=0" />
        <meta http-equiv="cache-control" content="no-cache" />
        <meta http-equiv="expires" content="0" />
        <meta http-equiv="expires" content="Tue, 01 Jan 1980 1:00:00 GMT" />
        <meta http-equiv="pragma" content="no-cache" />
        <link rel="stylesheet" type="text/css" href="{{ .Prefix }}.css" />
        {{- if .CustomerHeader }}{{ .CustomHeader }}{{ end }}
    <style>
        body, html {
            margin: 0;
            padding: 0;
        }

        .pycard {
            float: left;
            position: relative;
        }

        @media print {
            div {
                page-break-inside: avoid;
            }
        }
    </style>
    </head>
    <body>
        {{- range .RenderedCards }}
        <div class="pycard">
            {{ .Card }}
        </div>
        {{ endfor }}
    </body>
</html>
`
