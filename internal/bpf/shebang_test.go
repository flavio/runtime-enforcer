package bpf

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/rancher-sandbox/runtime-enforcer/internal/types/policymode"
	"github.com/stretchr/testify/require"
)

// createShebangScript creates a temporary executable script with the
// given shebang line and returns its path.
func createShebangScript(t *testing.T, interpreter string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "test.sh")
	require.NoError(t, os.WriteFile(path, []byte("#!"+interpreter+"\n"), 0755), "failed to write shebang script")
	return path
}

func TestShebangScript(t *testing.T) {
	runner, err := newCgroupRunner(t)
	require.NoError(t, err, "Failed to create cgroup runner")
	defer runner.close()

	const interpreter = "/usr/bin/true"
	scriptPath := createShebangScript(t, interpreter)

	// When a shebang script is executed, the LSM hook fires for
	// the script but not for the interpreter.
	require.NoError(t, runner.runAndFindCommand(&runCommandArgs{
		command:         scriptPath,
		channel:         learningChannel,
		shouldFindEvent: true,
	}), "script path must be learned via the LSM hook")

	// The interpreter itself must NOT appear as a learned event;
	// only the script path should be emitted.
	require.NoError(t, runner.runAndFindCommand(&runCommandArgs{
		command:         scriptPath,
		channel:         learningChannel,
		expectedPath:    "/usr/bin/true",
		shouldFindEvent: false,
	}), "interpreter must not be learned, only the script path")

	// Once a policy is active, learning events must stop being emitted.
	mockPolicyID := uint64(44)
	err = runner.manager.GetPolicyUpdateBinariesFunc()(
		mockPolicyID,
		[]string{scriptPath},
		AddValuesToPolicy,
	)
	require.NoError(t, err)

	err = runner.manager.GetPolicyModeUpdateFunc()(mockPolicyID, policymode.Protect, UpdateMode)
	require.NoError(t, err)

	err = runner.manager.GetCgroupPolicyUpdateFunc()(
		mockPolicyID, []uint64{runner.cgInfo.id}, AddPolicyToCgroups,
	)
	require.NoError(t, err)

	require.NoError(t, runner.runAndFindCommand(&runCommandArgs{
		command:         scriptPath,
		channel:         learningChannel,
		shouldFindEvent: false,
	}), "script path must not be learned once a policy is active")

	require.NoError(t, runner.runAndFindCommand(&runCommandArgs{
		command:         scriptPath,
		channel:         monitoringChannel,
		shouldFindEvent: false,
	}), "script must be allowed")

	require.NoError(t, runner.runAndFindCommand(&runCommandArgs{
		command:         scriptPath,
		channel:         monitoringChannel,
		expectedPath:    "/usr/bin/true",
		shouldFindEvent: false,
	}), "/usr/bin/true is not seen ebpf side")
}
