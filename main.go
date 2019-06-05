package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	rice "github.com/GeertJohan/go.rice"
	"github.com/fsnotify/fsnotify"
	"github.com/ghostsquad/go-timejumper"
	"github.com/go-kit/kit/log"
	"github.com/gorilla/websocket"
	oauth2ns "github.com/nmrshll/oauth2-noserver"
	"github.com/oklog/run"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

var ClientID string
var ClientSecret string

type CardData struct {
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	if conn == nil {
		fmt.Println("conn is nil")
	}
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

	jsbox := rice.MustFindBox("dist/js")
	liveReloadjs, err := jsbox.String("livereload.min.js")
	if err != nil {
		logger.Log("event", "exit", "reason", err)
		os.Exit(1)
	}

	if liveReloadjs == "" {
		logger.Log("event", "exit", "reason", "liveReloadJs is empty")
		os.Exit(1)
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

	fmt.Printf("%+v\n", ClientID)

	config := &oauth2.Config{
		ClientID:     ClientID,
		ClientSecret: ClientSecret,
		Endpoint:     google.Endpoint,
		Scopes:       []string{"https://www.googleapis.com/auth/spreadsheets.readonly"},
	}

	//authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	client, err := oauth2ns.AuthenticateUser(config)
	if err != nil {
		logger.Log("error", err)
		os.Exit(1)
	}

	srv, err := sheets.New(client.Client)
	if err != nil {
		logger.Log("error", fmt.Sprintf("Unable to retrieve Sheets client: %v", err))
		os.Exit(1)
	}

	err = watcher.Add("./content")

	if srv == nil {
		logger.Log("error", fmt.Sprintf("srv is nil: %v", err))
	}

	// for {
	// 	select {
	// 	case event := <-watcher.Events:
	// 		lastChangedSources := map[string]struct{}{event.Name: {}}

	// 		if !shouldRebuild(event.Name, event.Op) {
	// 			continue
	// 		}

	// 		for {
	// 			if len(lastChangedSources) < 1 {
	// 				break
	// 			}

	// 			rebuild <- lastChangedSources

	// 			// Zero out the last set of changes and start
	// 			// accumulating.
	// 			lastChangedSources = nil

	// 			// Wait until rebuild is finished. In the meantime,
	// 			// accumulate new events that come in on the watcher's
	// 			// channel and prepare for the next loop.
	// 		INNER_LOOP:
	// 			for {
	// 				select {
	// 				case <-rebuildDone:
	// 					break INNER_LOOP
	// 				case event := <-watchEvents:
	// 					if !shouldRebuild(event.Name, event.Op) {
	// 						continue
	// 					}

	// 					if lastChangedSources == nil {
	// 						lastChangedSources = make(map[string]struct{})
	// 					}

	// 					lastChangedSources[event.Name] = struct{}{}
	// 				}
	// 			}
	// 		}

	// 		logger.Log("event", event)
	// 	}
	// }

	// t := template.Must(template.New("email.tmpl").Parse(indexTmpl))

	// err := t.Execute(os.Stdout, getCardData())
	// if err != nil {
	//     logger.Log("event", "exit", "reason", errors.Wrap(err, "Failed to render template!")
	// 	os.Exit(1)
	// }

	logger.Log("event", "exit", "reason", runGroup.Run())
}

func getCardData() *CardData {
	return &CardData{}
}

func shouldRebuild(path string, op fsnotify.Op) bool {
	base := filepath.Base(path)

	// Mac OS' worst mistake.
	if base == ".DS_Store" {
		return false
	}

	// Vim creates this temporary file to see whether it can write
	// into a target directory. It screws up our watching algorithm,
	// so ignore it.
	if base == "4913" {
		return false
	}

	// A special case, but ignore creates on files that look like
	// Vim backups.
	if strings.HasSuffix(base, "~") {
		return false
	}

	return true
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
