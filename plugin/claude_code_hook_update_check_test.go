package plugin_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestClaudeHookHelpersDisableUpdateCheck pins the shell half of the hook
// update-check fix: sourcing _helpers.sh must export ENGRAM_NO_UPDATE_CHECK=1
// so every binary call a hook makes (instance-id, sync --import, protocol-mode)
// skips the GitHub release check. `sync --import` is not excluded in Go, so this
// export is the only thing that keeps it off the hook's time budget. printenv
// reports exported variables only, so a plain assignment would fail here.
func TestClaudeHookHelpersDisableUpdateCheck(t *testing.T) {
	requireHookBinaries(t)
	helper := filepath.Join(repoRoot(t), "plugin", "claude-code", "scripts", "_helpers.sh")
	cmd := exec.Command("bash", "-c", `source "$1"; printenv ENGRAM_NO_UPDATE_CHECK`, "bash", bashScriptPath(t, helper))
	cmd.Env = make([]string, 0, len(os.Environ()))
	for _, entry := range os.Environ() {
		if !strings.HasPrefix(strings.ToUpper(entry), "ENGRAM_NO_UPDATE_CHECK=") {
			cmd.Env = append(cmd.Env, entry)
		}
	}
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("_helpers.sh must export ENGRAM_NO_UPDATE_CHECK: %v (output %q)", err, output)
	}
	if got := strings.TrimSpace(string(output)); got != "1" {
		t.Fatalf("ENGRAM_NO_UPDATE_CHECK after sourcing _helpers.sh = %q, want %q", got, "1")
	}
}
