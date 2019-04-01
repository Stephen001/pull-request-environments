package main

import (
    "fmt"
    "net/http"
    "gopkg.in/go-playground/webhooks.v5/github"
)

const (
    path = "/webhooks"
)

func main() {
    hook, _ := github.New(github.Options.Secret("hello"))

    http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
        payload, err := hook.Parse(r, github.PullRequestEvent)
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
	http.ListenAndServe(":3000", nil)
}
