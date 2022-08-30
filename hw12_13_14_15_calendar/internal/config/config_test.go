package config

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const tmpConfigFileName = "tmp.toml"

func TestConfig(t *testing.T) {
	testCases := []struct {
		name     string
		expected Config
	}{
		{
			name: "psql full",
			expected: Config{
				Logger: LoggerConf{
					Level:   "INFO",
					OutPath: "output.log",
				},
				Storage: StorageConf{
					StorageType: "psql",
					User:        "user",
					Password:    "pass",
					Host:        "host.com",
					Port:        "6543",
					DBName:      "db",
				},
				HTTPServer: ServerConf{
					Host: "0.0.0.0",
					Port: "1234",
				},
				GRPCServer: ServerConf{
					Host: "0.0.0.0",
					Port: "4321",
				},
			},
		},
		{
			name: "inmemory full",
			expected: Config{
				Logger: LoggerConf{
					Level:   "ERROR",
					OutPath: "stdout",
				},
				Storage: StorageConf{
					StorageType: "inmemory",
				},
				HTTPServer: ServerConf{
					Host: "0.0.0.0",
					Port: "1234",
				},
				GRPCServer: ServerConf{
					Host: "0.0.0.0",
					Port: "4321",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			te := setUpTestEnv(t, getConfigString(tc.expected))
			defer te.tearDown(t)

			got, err := NewConfig(te.configFile.Name())
			require.NoError(t, err)

			require.Equal(t, tc.expected, got)
		})
	}

	t.Run("default", func(t *testing.T) {
		expected := Config{
			Logger: LoggerConf{
				Level:   DefaultLogLevel,
				OutPath: DefaultLogOutFile,
			},
			Storage: StorageConf{
				StorageType: InmemoryStorageType,
			},
			HTTPServer: ServerConf{
				Host: DefaultHTTPHost,
				Port: DefaultHTTPPort,
			},
			GRPCServer: ServerConf{
				Host: DefaultGRPCHost,
				Port: DefaultGRPCPort,
			},
			AMQP: AMQPConf{
				Uri:   DefaultAMQPUri,
				Queue: DefaultAMQPQueue,
			},
			Scheduler: SchedulerConf{
				RemindBefore: time.Hour * 1,
			},
		}
		te := setUpTestEnv(t, ``)
		defer te.tearDown(t)

		got, err := NewConfig(te.configFile.Name())
		require.NoError(t, err)

		require.Equal(t, expected, got)
	})
}

func TestConfigNoFile(t *testing.T) {
	t.Run("no such file", func(t *testing.T) {
		_, err := NewConfig("invalidfile")
		require.Error(t, err)
	})
}

func TestConfigNoFieldInPSQLStorage(t *testing.T) {
	testCases := []struct {
		name, fileContent string
	}{
		{
			name: "no user",
			fileContent: `[storage]
			storate_type = "psql"`,
		},
		{
			name: "no pass",
			fileContent: `[storage]
			storate_type = "psql"
			user = "postgres"`,
		},
		{
			name: "no db",
			fileContent: `[storage]
			storate_type = "psql"
			user = "postgres"
			password = "password"`,
		},
		{
			name: "no host",
			fileContent: `[storage]
			storate_type = "psql"
			user = "postgres"
			password = "password"
			db = "db"`,
		},
		{
			name: "no port",
			fileContent: `[storage]
			storate_type = "psql"
			user = "postgres"
			password = "password"
			db = "db"
			host = "host"`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			te := setUpTestEnv(t, tc.fileContent)
			defer te.tearDown(t)

			_, err := NewConfig(te.configFile.Name())
			require.Error(t, err)
		})
	}
}

type testEnv struct {
	configFile *os.File
	sourceText string
}

func (te *testEnv) tearDown(t *testing.T) {
	t.Helper()
	err := te.configFile.Close()
	if err != nil {
		t.Fatal(err)
	}
	err = os.Remove(te.configFile.Name())
	if err != nil {
		t.Fatal(err)
	}
}

func setUpTestEnv(t *testing.T, sourceText string) testEnv {
	t.Helper()
	tmpSourceFile, err := os.CreateTemp("", tmpConfigFileName)
	if err != nil {
		t.Fatal(err)
	}
	_, err = tmpSourceFile.Write([]byte(sourceText))
	if err != nil {
		t.Fatal(err)
	}

	return testEnv{
		configFile: tmpSourceFile,
		sourceText: sourceText,
	}
}

func getConfigString(c Config) string {
	return fmt.Sprintf(`[logger]
level = "%s"
file = "%s"

[storage]
storage_type = "%s"
user = "%s"
password = "%s"
db = "%s"
host = "%s"
port = "%s"

[http]
host = "%s"
port = "%s"

[grpc]
host = "%s"
port = "%s"

[amqp]
uri = "%s"
queue = "%s"

[scheduler]
remind_before = "%s"
`,
		c.Logger.Level,
		c.Logger.OutPath,
		c.Storage.StorageType,
		c.Storage.User,
		c.Storage.Password,
		c.Storage.DBName,
		c.Storage.Host,
		c.Storage.Port,
		c.HTTPServer.Host,
		c.HTTPServer.Port,
		c.GRPCServer.Host,
		c.GRPCServer.Port,
		c.AMQP.Uri,
		c.AMQP.Queue,
		c.Scheduler.RemindBefore,
	)
}
