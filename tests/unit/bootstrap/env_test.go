package bootstrap_test

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/heaveless/dbz-api/internal/bootstrap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewEnv_LoadsConfigFromDotEnv(t *testing.T) {
	tempDir := t.TempDir()

	origWD, err := os.Getwd()
	require.NoError(t, err)

	require.NoError(t, os.Chdir(tempDir))
	t.Cleanup(func() {
		_ = os.Chdir(origWD)
	})

	envContent := []byte(`
APP_ENV=development
APP_PORT=8080
DB_HOST=localhost
DB_PORT=27017
DB_NAME=dbz
API_URI=https://example.com
`)
	err = os.WriteFile(filepath.Join(tempDir, ".env"), envContent, 0o644)
	require.NoError(t, err)

	var buf bytes.Buffer
	origOutput := log.Writer()
	log.SetOutput(&buf)
	t.Cleanup(func() {
		log.SetOutput(origOutput)
	})

	env := bootstrap.NewEnv()

	require.NotNil(t, env)
	assert.Equal(t, "development", env.AppEnv)
	assert.Equal(t, "8080", env.AppPort)
	assert.Equal(t, "localhost", env.DBHost)
	assert.Equal(t, "27017", env.DBPort)
	assert.Equal(t, "dbz", env.DBName)
	assert.Equal(t, "https://example.com", env.ApiUri)

	logOutput := buf.String()
	assert.Contains(t, logOutput, "The App is running in development env")
}
