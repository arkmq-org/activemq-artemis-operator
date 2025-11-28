---
title: "Debugging Go Code"
subtitle: "A comprehensive guide to debugging the operator"
description: "Learn how to set up debugging sessions using Delve, VS Code, Cursor, and Neovim"
type: "tutorial"
draft: false
---

# Debugging Go

This tutorial walks you through setting up a debugging session for the ActiveMQ Artemis Operator. We use [Delve](https://github.com/go-delve/delve), the standard debugger for Go.

## Prerequisites

### Install Delve

```bash
go install github.com/go-delve/delve/cmd/dlv@latest
```

### Set Up the Operator for Local Debugging

Install the CRDs on your cluster before debugging:

```bash
# Install CRDs on your cluster
make install

# When done, clean up
make uninstall
```

---

## Part 1: CLI Debugging with Delve

This method is essential for debugging on remote servers or containers where you don't have a GUI.

### Start an Interactive Session

For the operator, use:

```bash
make debug DEFAULT_OPERATOR_NAMESPACE=my-namespace  # Sets DEFAULT_OPERATOR_NAMESPACE
```

This compiles with optimizations disabled and drops you into the `(dlv)` interactive shell.

### Setting Breakpoints

Break by file and line:

```
break controllers/activemqartemis_controller.go:150
```

Or by function name:

```
break controllers.(*ActiveMQArtemisReconciler).Reconcile
```

### Running and Inspecting

| Command | Description |
|---------|-------------|
| `continue` (c) | Run until next breakpoint |
| `next` (n) | Step to next line (don't enter functions) |
| `step` (s) | Step into function call |
| `stepout` | Step out of current function |
| `locals` | Show local variables |
| `args` | Show function arguments |
| `print <expr>` | Print a variable or expression |
| `quit` | Exit the debugger |

### Example Session

```
(dlv) break controllers/activemqartemis_controller.go:200
Breakpoint 1 set at 0x... for controllers.(*ActiveMQArtemisReconciler).Reconcile()
(dlv) continue
> controllers.(*ActiveMQArtemisReconciler).Reconcile() ./controllers/activemqartemis_controller.go:200
(dlv) locals
request = {NamespacedName: {Namespace: "default", Name: "my-broker"}}
(dlv) print request.Name
"my-broker"
(dlv) next
(dlv) quit
```

---

## Part 2: Enhanced Debugging (VS Code, Cursor, Neovim)

All three editors can use the same `.vscode/launch.json` configuration file.

> [!Note]
> Neovim with `nvim-dap` reads `.vscode/launch.json` automatically, so you only need one config file for all editors!

### Configuration: launch.json

Create a `.vscode/launch.json` with these configurations:

```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Debug Operator (launch)",
            "type": "go",
            "request": "launch",
            "mode": "debug",
            "program": "${workspaceFolder}",
            "env": {
                "DEFAULT_OPERATOR_NAMESPACE": "activemq-artemis-operator",
                "WATCH_NAMESPACE": "",
                "POD_NAME": "local-operator-debug",
                "OPERATOR_NAME": "activemq-artemis-operator"
            },
            "args": [
                "--leader-elect=false"
            ]
        },
        {
            "name": "Attach to Delve (make debug-remote)",
            "type": "go",
            "request": "attach",
            "mode": "remote",
            "remotePath": "${workspaceFolder}",
            "port": 2345,
            "host": "127.0.0.1"
        }
    ]
}
```

---

## Part 2a: VS Code / Cursor

### Step 1: Install the Go Extension

Ensure you have the official **Go** extension by the Go Team at Google installed.

### Step 2: Set Breakpoints

1. Open a source file (e.g., `controllers/activemqartemis_controller.go`)
2. Click to the left of a line number
3. A red dot appears — that's your breakpoint

### Step 4: Choose a Debug Configuration

To select a debug configuration:

1. Click the **Run and Debug** icon in the left sidebar (play button with a bug)
2. At the top of the panel, click the **dropdown menu**
3. Select from the available configurations:
   - **Debug Operator (launch)** — Default. IDE launches everything.
   - **Attach to Delve (make debug-remote)** — Connects to an already-running Delve server.

### Step 5: Start Debugging

**Option A: Direct Launch (recommended)**

1. Select **"Debug Operator (launch)"** from the dropdown
2. Press **F5** or click the green Play button
3. The program starts and pauses at your breakpoint

This is the simplest approach — the IDE handles everything.

**Option B: Attach to Running Delve (advanced)**

This is optional — only use when you need custom setup or want to reattach without restarting:

1. In a terminal, start the Delve server:
   ```bash
   make debug-remote DEFAULT_OPERATOR_NAMESPACE=my-namespace  # Sets DEFAULT_OPERATOR_NAMESPACE
   ```
2. Select **"Attach to Delve (make debug-remote)"** from the dropdown
3. Press **F5** to connect

> This option works in both VS Code/Cursor and Neovim.

### Debug Controls

| Button | Shortcut | Action |
|--------|----------|--------|
| ▶️ Continue | F5 | Run to next breakpoint |
| ⏭️ Step Over | F10 | Next line (skip function internals) |
| ⬇️ Step Into | F11 | Enter function call |
| ⬆️ Step Out | Shift+F11 | Exit current function |
| 🔄 Restart | Ctrl+Shift+F5 | Restart session |
| ⏹️ Stop | Shift+F5 | End session |

### Debug Sidebar Panels

- **Variables**: Shows locals and globals. Expand objects to see fields.
- **Watch**: Add expressions to evaluate live (e.g., `len(pods.Items)`)
- **Call Stack**: See the full call chain. Click to jump between frames.
- **Breakpoints**: Manage all breakpoints, enable/disable, add conditions.

---

## Part 2b: Neovim with LazyVim

LazyVim has built-in support for Go debugging via `nvim-dap-go`. It reads the same `.vscode/launch.json` file, so no extra configuration is needed!

### Step 1: Enable Required Extras

Open Neovim and run `:LazyExtras`, then enable:

- [x] `lang.go` — Go language support
- [x] `dap.core` — Debug Adapter Protocol (keybindings and UI)

Restart Neovim after enabling.

> ⚠️ If you skip `dap.core`, you'll have the adapter but no keybindings!

### Step 2: Verify Delve is Installed

Run `:Mason` and ensure `delve` appears as installed.

### Step 3: Set Breakpoints

1. Open a source file
2. Move cursor to the target line
3. Press `<leader>db` (Space → d → b)
4. A dot appears in the gutter

### Step 4: Start Debugging

1. Press `<leader>dc` (Space → d → c)
2. Select a configuration from the list (these come from `.vscode/launch.json`):
   - **Debug Operator (launch)** — Recommended
   - **Debug Operator (single namespace)**
   - **Attach to Delve (make debug-remote)**

### Keybindings

| Action | Shortcut | Description |
|--------|----------|-------------|
| Continue/Start | `<leader>dc` | Start or resume execution |
| Step Over | `<leader>dO` | Next line (don't enter functions) |
| Step Into | `<leader>di` | Enter function call |
| Step Out | `<leader>do` | Exit current function |
| Toggle Breakpoint | `<leader>db` | Add/remove breakpoint |
| Terminate | `<leader>dt` | Stop the debug session |
| Toggle UI | `<leader>du` | Open/close debug windows |
| Eval | `<leader>de` | Evaluate expression under cursor |

### Troubleshooting

**Missing keybindings?** Ensure `dap.core` is enabled in `:LazyExtras`.

**Delve not found?** Run `:Mason` and install `delve`.

---

## Part 3: Debugging Tests

The operator uses Ginkgo for integration/E2E tests.

### Test Configuration in launch.json

Add this object to your `launch.json` file:

```json
{
    "name": "Debug Test (current file)",
    "type": "go",
    "request": "launch",
    "mode": "test",
    "program": "${fileDirname}",
    "env": {
        "USE_EXISTING_CLUSTER": "true",
        "DEFAULT_OPERATOR_NAMESPACE": "activemq-artemis-operator",
        "WATCH_NAMESPACE": "",
        "KUBECONFIG": "${env:HOME}/.kube/config"
    },
    "args": [
        "-ginkgo.v",
        "-ginkgo.fail-fast"
    ]
}
```

> **Note:** `${fileDirname}` means "directory of the currently open file". Open any `_test.go` file and this config will run tests in that package.

### Debugging a Single Ginkgo Test (Recommended Workflow)

#### Step 1: Focus the test with `FIt`

Find the test you want to debug and change `It` to `FIt`:

```go
// Before:
It("configurationManaged default should be true", func() {
    // test code...
})

// After:
FIt("configurationManaged default should be true", func() {
    // test code...
})
```

Ginkgo will skip all other tests and only run the focused one.

#### Step 2: Set breakpoints

Click in the gutter (left of line numbers) where you want to pause execution.

#### Step 3: Start debugging

**VS Code / Cursor:**
1. Open the test file you want to debug
2. Open **Run and Debug** sidebar
3. Select **"Debug Test (current file)"**
4. Press **F5**

**Neovim:**
1. Open the test file you want to debug
2. Press `<leader>dc`
3. Select **"Debug Test (current file)"**

#### Step 4: Debug!

The test runs, stops at your breakpoint. Inspect variables, step through code.

#### Step 5: Clean up

> ⚠️ **Don't forget** to change `FIt` back to `It` before committing!

```go
It("configurationManaged default should be true", func() {
```

### Alternative: Filter by Regex (no code changes)

If you don't want to modify source code, use `-ginkgo.focus` in launch.json:

```json
"args": [
    "-ginkgo.v",
    "-ginkgo.focus", "configurationManaged default"
]
```

The focus accepts a regex pattern. Edit this string to match different tests.

### Alternative: Filter by Labels

Some tests have labels like `Label("slow")`. Filter them:

```json
"args": [
    "-ginkgo.v",
    "-ginkgo.label-filter", "queue-config-defaults"
]
```

Or exclude labels:

```json
"args": [
    "-ginkgo.v",
    "-ginkgo.label-filter", "!slow && !do"
]
```

### Standard Go Tests

For simple `func TestXxx(t *testing.T)` tests (non-Ginkgo):

1. Open the `_test.go` file
2. Set a breakpoint inside the test function
3. **VS Code/Cursor**: Click "debug test" above the function (CodeLens)
4. **Neovim**: `<leader>dc` → select "Debug test"

---

## Makefile Targets Reference

| Target | Description |
|--------|-------------|
| `make install` | Install CRDs on cluster |
| `make uninstall` | Uninstall CRDs from cluster |
| `make debug` | Start Delve in interactive CLI mode |
| `make debug-remote` | Start Delve on port 2345 for IDE attach |

### Examples

```bash
# Debug watching all namespaces
# Install CRDs first
make install

# Debug (uses DEFAULT_OPERATOR_NAMESPACE Makefile variable, sets DEFAULT_OPERATOR_NAMESPACE env)
make debug DEFAULT_OPERATOR_NAMESPACE=activemq-artemis-operator

# Debug watching a single namespace
make debug DEFAULT_OPERATOR_NAMESPACE=my-ns WATCH_NAMESPACE=my-ns

# Start remote debugger for IDE
make debug-remote DEFAULT_OPERATOR_NAMESPACE=my-ns
```
---

## See Also

- [Delve Documentation](https://github.com/go-delve/delve/tree/master/Documentation)
- [VS Code Go Debugging](https://code.visualstudio.com/docs/languages/go#_debugging)
- [nvim-dap-go](https://github.com/leoluz/nvim-dap-go)
- [Ginkgo Filtering](https://onsi.github.io/ginkgo/#filtering-specs)
