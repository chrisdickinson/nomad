//go:build ent
// +build ent

package command

import (
	"strings"
	"testing"

	"github.com/hashicorp/nomad/ci"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
	"github.com/shoenig/test/must"
)

func TestQuotaInspectCommand_Implements(t *testing.T) {
	ci.Parallel(t)
	var _ cli.Command = &QuotaInspectCommand{}
}

func TestQuotaInspectCommand_Fails(t *testing.T) {
	ci.Parallel(t)
	ui := cli.NewMockUi()
	cmd := &QuotaInspectCommand{Meta: Meta{Ui: ui}}

	// Fails on misuse
	code = cmd.Run([]string{"some", "bad", "args"})
	must.One(t, code)

	must.StrContains(t, ui.ErrorWriter.String(), commandErrorText(cmd))
	ui.ErrorWriter.Reset()

	code = cmd.Run([]string{"-address=nope", "foo"})
	must.One(t, code)

	must.StrContains(t, ui.ErrorWriter.String(), "retrieving quota")
	ui.ErrorWriter.Reset()
}

func TestQuotaInspectCommand_Run(t *testing.T) {
	ci.Parallel(t)

	// Create a server
	srv, client, url := testServer(t, true, nil)
	defer srv.Shutdown()

	ui := cli.NewMockUi()
	cmd := &QuotaInspectCommand{Meta: Meta{Ui: ui}}

	// Create a quota to delete
	qs := testQuotaSpec()
	_, err := client.Quotas().Register(qs, nil)
	must.NoError(t, err)

	// Delete a quota
	code := cmd.Run([]string{"-address=" + url, qs.Name})
	must.Zero(t, code)

	out := ui.OutputWriter.String()
	if !strings.Contains(out, "Usages") || !strings.Contains(out, qs.Name) {
		t.Fatalf("expected quota, got: %s", out)
	}

	// List json
	must.Zero(t, cmd.Run([]string{"-address=" + url, "-json", allocID}))

	outJson := *api.Quotas{}
	err = json.Unmarshal(ui.OutputWriter.Bytes(), &outJson)
	must.NoError(t, err)

	ui.OutputWriter.Reset()

	// Go template to format the output
	code = cmd.Run([]string{"-address=" + url, "-t", "{{ .Name }}", allocID})
	must.Zero(t, code)

	out = ui.OutputWriter.String()
	must.StrContains(t, out, "test-quota")

	ui.OutputWriter.Reset()
}

func TestQuotaInspectCommand_AutocompleteArgs(t *testing.T) {
	ci.Parallel(t)
	assert := assert.New(t)

	srv, client, url := testServer(t, true, nil)
	defer srv.Shutdown()

	ui := cli.NewMockUi()
	cmd := &QuotaInspectCommand{Meta: Meta{Ui: ui, flagAddress: url}}

	// Create a quota
	qs := testQuotaSpec()
	_, err := client.Quotas().Register(qs, nil)
	must.NoError(t, err)

	args := complete.Args{Last: "t"}
	predictor := cmd.AutocompleteArgs()

	res := predictor.Predict(args)
	must.One(t, len(res))
	must.StrContains(t, qs.Name, res[0])
}
