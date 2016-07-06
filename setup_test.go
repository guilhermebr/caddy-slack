package slack

import (
	"testing"

	"github.com/mholt/caddy"
	"github.com/mholt/caddy/caddyhttp/httpserver"
)

func TestSetup(t *testing.T) {
	c := caddy.NewTestController("http", `slack T024Z4TDT/B1NNW1YBZ/R08d98T94dPfsN5adHABlSs3`)
	err := setup(c)

	if err != nil {
		t.Errorf("Expected no errors, got: %v", err)
	}
	cfg := httpserver.GetConfig(c)
	mids := cfg.Middleware()
	if mids == nil {
		t.Fatal("Expected middleware, was nil instead")
	}

	handler := mids[0](httpserver.EmptyNext)
	myHandler, ok := handler.(Slack)

	if !ok {
		t.Fatalf("Expected handler to be type Slack, got: %#v", handler)
	}

	if myHandler.Rules[0].Token != "T024Z4TDT/B1NNW1YBZ/R08d98T94dPfsN5adHABlSs3" {
		t.Errorf("Wrong token: %+v", myHandler.Rules[0].Token)
	}

}
