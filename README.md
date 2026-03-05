# cli-wizard-core

[![CI](https://github.com/Bibi40k/cli-wizard-core/actions/workflows/ci.yml/badge.svg)](https://github.com/Bibi40k/cli-wizard-core/actions/workflows/ci.yml)
[![Release](https://github.com/Bibi40k/cli-wizard-core/actions/workflows/release.yml/badge.svg)](https://github.com/Bibi40k/cli-wizard-core/actions/workflows/release.yml)
[![Go Version](https://img.shields.io/github/go-mod/go-version/bibi40k/cli-wizard-core)](https://github.com/Bibi40k/cli-wizard-core/blob/main/go.mod)
[![License](https://img.shields.io/github/license/bibi40k/cli-wizard-core)](https://github.com/Bibi40k/cli-wizard-core/blob/main/LICENSE)

Reusable core primitives for interactive CLI wizards in Go.

## What it provides

- `Session`: draft lifecycle orchestration (load/start/stop/finalize)
- `RunSteps`: deterministic step runner with start/done callbacks
- `Manage`: generic create/edit/delete/default action flow for list-based resources
- `FormatMenuLabel`: aligned two-column menu labels (`[tag]` + text)
- `Colorize`: ANSI color wrapper for labels/messages
- `NewUserError` / `WithHint`: attach actionable hints to user-facing errors
- `FormatCLIError`: consistent colored error+hint output for CLIs
- `ReconcileWithTemplate`: generic template sync for config objects (`added/removed/missing-required` report)

The package is intentionally transport-agnostic: no direct dependency on survey/readline/TTY, no provider-specific logic (VMware, Talos, etc.).

## Install

```bash
go get github.com/Bibi40k/cli-wizard-core
```

## Usage

```go
import wizard "github.com/Bibi40k/cli-wizard-core"

session := wizard.NewSession(
    targetPath,
    draftPath,
    &state,
    isEmpty,
    loadDraftFn,
    startDraftFn,
    finalizeFn,
)

_ = wizard.RunSteps([]wizard.Step{
    {Name: "Step 1", Run: step1},
    {Name: "Step 2", Run: step2},
}, onStart, onDone)

template := map[string]any{
    "vm": map[string]any{
        "name":       "",
        "profile":    "talos",
        "ip_address": "",
    },
}
current := map[string]any{
    "vm": map[string]any{
        "name":     "cp-01",
        "username": "",
    },
}

out, report, err := wizard.ReconcileWithTemplate(current, template, wizard.ReconcileOptions{
    DropUnknown:   true,
    RequiredPaths: []string{"vm.name", "vm.ip_address"},
})
_ = out
_ = report
_ = err
```
