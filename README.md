# cli-wizard-core

Reusable core primitives for interactive CLI wizards in Go.

## What it provides

- `Session`: draft lifecycle orchestration (load/start/stop/finalize)
- `RunSteps`: deterministic step runner with start/done callbacks

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
```
