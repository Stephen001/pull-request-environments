package main

import (
    "os"
    "fmt"
    "path"
    "net/http"
    "flag"
    "io/ioutil"
    "github.com/fsnotify/fsnotify"
    "gopkg.in/go-playground/webhooks.v5/github"
)

const (
    uriPath = "/"
)

var (
    phrasePath = "/etc/secrets/webhook/phrase"
)

func init() {
    flag.StringVar(&phrasePath, "phrase", "/etc/secrets/webhook/phrase", "The path to a file containing the phrase to secure webhook calls with.")
}

func readPhrase(phrasePath string) (string, error) {
    phrase, err := ioutil.ReadFile(phrasePath)
    if err != nil {
        fmt.Println("Could not read webhook phrase file:", err)
    }
    return string(phrase), err
}

func main() {
    flag.Parse()

    var hook * github.Webhook = nil

    phrase, err := readPhrase(phrasePath)
    if err != nil {
        os.Exit(1)
    }

    hook, _ = github.New(github.Options.Secret(phrase))
    fmt.Println("Webhook phrase loaded successfully.")

    watcher, err := fsnotify.NewWatcher()
    if err != nil {
        fmt.Println("Could not create file watcher:", err)
        os.Exit(2)
    }
    defer watcher.Close()

    go func() {
        for {
            select {
            case event, ok := <-watcher.Events:
                if !ok {
                    return
                }
                if event.Op&fsnotify.Write == fsnotify.Write {
                    if event.Name == path.Base(phrasePath) {
                        phrase, err := readPhrase(phrasePath)
                        if err != nil {
                            fmt.Println("Not updating webhook phrase:", err)
                        } else {
                            hook, _ = github.New(github.Options.Secret(phrase))
                            fmt.Println("Webhook phrase has been updated.")
                        }
                    }
                }
            case err, ok := <-watcher.Errors:
                if !ok {
                    return
                }
                fmt.Println("File watcher error:", err)
            }
        }
    }()

    err = watcher.Add(phrasePath)
    if err != nil {
        fmt.Println("Could not watch phrase file:", err)
        os.Exit(2)
    }

    http.HandleFunc(uriPath, func(w http.ResponseWriter, r *http.Request) {
        payload, err := (*hook).Parse(r, github.PullRequestEvent)
        if err != nil {
            if err == github.ErrEventNotFound {
                // ok event wasn't one of the ones asked to be parsed
            }
        }

        switch payload.(type) {

        case github.PullRequestPayload:
            pullRequest := payload.(github.PullRequestPayload)
            // Do whatever you want from here...
            fmt.Printf("%+v", pullRequest)
        }
    })
    fmt.Println("Listening on port 3000")
    http.ListenAndServe(":3000", nil)
}
