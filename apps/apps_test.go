package apps

import (
	"testing"

	"github.com/gjhenrique/yafl/internal/test"
	"github.com/stretchr/testify/require"
)

func TestEmptyDirectory(t *testing.T) {
	workspace := test.SetupWorkspace(t)
	defer workspace.RemoveWorkspace()

	t.Setenv(customDesktopsEnv, workspace.Dir)
	names, err := FormattedApplicationNames()
	require.NoError(t, err)

	require.Len(t, names, 0)

}

// TODO: Write tests after refactoring modes structure
