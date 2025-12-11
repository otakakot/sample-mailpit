package main_test

import (
	"crypto/rand"
	"fmt"
	"strings"
	"testing"

	"github.com/axllent/mailpit/sendmail/cmd"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"

	"github.com/otakakot/sample-mailpit/pkg/mailpit/client"
	"github.com/otakakot/sample-mailpit/pkg/mailpit/client/messages"
)

const (
	smtpHost = "localhost:1025"
	httpHost = "localhost:8025"
)

func TestMailpit(t *testing.T) {
	from := "noreply@example.com"

	to := rand.Text() + "@example.com"

	subject := "DevelopTest"

	body := "From Develop Test"

	msg := []byte(strings.ReplaceAll(fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\n\n%s", from, strings.Join([]string{to}, ","), subject, body), "\n", "\r\n"))

	if err := cmd.Send(smtpHost, from, []string{to}, msg); err != nil {
		t.Fatalf("Failed to send email: %v", err)
	}

	transport := httptransport.New(httpHost, "", []string{"http"})
	cli := client.New(transport, strfmt.Default)

	res, err := cli.Messages.SearchParams(messages.NewSearchParamsParams().WithQuery(to))
	if err != nil {
		t.Fatalf("Failed to search messages: %v", err)
	}

	for _, msg := range res.Payload.Messages {
		if msg.From.Address != from {
			t.Errorf("Expected sender %q, got %q", from, msg.From.Address)
		}
		if msg.To[0].Address != to {
			t.Errorf("Expected recipient %q, got %q", to, msg.To[0].Address)
		}
		if msg.Subject != subject {
			t.Errorf("Expected subject %q, got %q", subject, msg.Subject)
		}
		if msg.Snippet != body {
			t.Errorf("Expected body snippet %q, got %q", body, msg.Snippet)
		}
	}
}
