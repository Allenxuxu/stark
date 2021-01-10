package env

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/Allenxuxu/stark/config/source"
)

func TestEnv_Read(t *testing.T) {
	expected := map[string]map[string]string{
		"database": {
			"host":       "localhost",
			"password":   "password",
			"datasource": "user:password@tcp(localhost:port)/db?charset=utf8mb4&parseTime=True&loc=Local",
		},
	}

	os.Setenv("DATABASE_HOST", "localhost")
	os.Setenv("DATABASE_PASSWORD", "password")
	os.Setenv("DATABASE_DATASOURCE", "user:password@tcp(localhost:port)/db?charset=utf8mb4&parseTime=True&loc=Local")

	source := NewSource()
	c, err := source.Read()
	if err != nil {
		t.Error(err)
	}

	var actual map[string]interface{}
	if err := json.Unmarshal(c.Data, &actual); err != nil {
		t.Error(err)
	}

	actualDB := actual["database"].(map[string]interface{})

	for k, v := range expected["database"] {
		a := actualDB[k]

		if a != v {
			t.Errorf("expected %v got %v", v, a)
		}
	}
}

func TestEnvvar_Prefixes(t *testing.T) {
	os.Setenv("APP_DATABASE_HOST", "localhost")
	os.Setenv("APP_DATABASE_PASSWORD", "password")
	os.Setenv("VAULT_ADDR", "vault:1337")
	os.Setenv("STARK_REGISTRY", "mdns")

	var prefixtests = []struct {
		prefixOpts   []source.Option
		expectedKeys []string
	}{
		{[]source.Option{WithPrefix("APP", "STARK")}, []string{"app", "stark"}},
		{[]source.Option{WithPrefix("STARK"), WithStrippedPrefix("APP")}, []string{"database", "stark"}},
		{[]source.Option{WithPrefix("STARK"), WithStrippedPrefix("APP")}, []string{"database", "stark"}},
	}

	for _, pt := range prefixtests {
		source := NewSource(pt.prefixOpts...)

		c, err := source.Read()
		if err != nil {
			t.Error(err)
		}

		var actual map[string]interface{}
		if err := json.Unmarshal(c.Data, &actual); err != nil {
			t.Error(err)
		}

		// assert other prefixes ignored
		if l := len(actual); l != len(pt.expectedKeys) {
			t.Errorf("expected %v top keys, got %v", len(pt.expectedKeys), l)
		}

		for _, k := range pt.expectedKeys {
			if !containsKey(actual, k) {
				t.Errorf("expected key %v, not found", k)
			}
		}
	}
}

func TestEnvvar_WatchNextNoOpsUntilStop(t *testing.T) {
	src := NewSource(WithStrippedPrefix("STARK_"))
	w, err := src.Watch()
	if err != nil {
		t.Error(err)
	}

	go func() {
		time.Sleep(50 * time.Millisecond)
		_ = w.Stop()
	}()

	if _, err := w.Next(); err != source.ErrWatcherStopped {
		t.Errorf("expected watcher stopped error, got %v", err)
	}
}

func containsKey(m map[string]interface{}, s string) bool {
	for k := range m {
		if k == s {
			return true
		}
	}
	return false
}
