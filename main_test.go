package main_test

import (
	"crypto/rand"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/axllent/mailpit/sendmail/cmd"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/otakakot/sample-mailpit/pkg/mailpit/client"
	"github.com/otakakot/sample-mailpit/pkg/mailpit/client/messages"
	"github.com/resend/resend-go/v3"
)

const (
	Domain   = "example.com"
	SMTPAddr = Domain + ":1025"
	HTTPAddr = "https://example.com"
)

func TestLocal(t *testing.T) {
	from := "noreply@example.com"

	to := rand.Text() + "@" + Domain

	subject := "DevelopTest"

	body := "From Develop Test"

	msg := []byte(strings.ReplaceAll(fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\n\n%s", from, strings.Join([]string{to}, ","), subject, body), "\n", "\r\n"))

	if err := cmd.Send("localhost:1025", from, []string{to}, msg); err != nil {
		t.Fatalf("Failed to send email: %v", err)
	}

	transport := httptransport.New("localhost:8025", "", []string{"http"})

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

func TestDevelop(t *testing.T) {
	from := "noreply@example.com"

	to := rand.Text() + "@" + Domain

	subject := "DevelopTest"

	body := "From Develop Test"

	msg := []byte(strings.ReplaceAll(fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\n\n%s", from, strings.Join([]string{to}, ","), subject, body), "\n", "\r\n"))

	if err := cmd.Send(SMTPAddr, from, []string{to}, msg); err != nil {
		t.Fatalf("Failed to send email: %v", err)
	}

	transport := httptransport.New(Domain, "", []string{"https"})

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

func TestStaging(t *testing.T) {
	key := ""

	resendCli := resend.NewClient(key)

	from := "noreply@example.com"

	to := rand.Text() + "@" + Domain

	subject := "StagingTest"

	body := "From Staging Test"

	params := &resend.SendEmailRequest{
		From:    from,
		To:      []string{to},
		Subject: subject,
		Text:    body,
	}

	sent, err := resendCli.Emails.Send(params)
	if err != nil {
		t.Fatalf("Failed to send email: %v", err)
	}
	t.Log(sent.Id)

	time.Sleep(10 * time.Second)

	transport := httptransport.New(Domain, "", []string{"https"})

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
