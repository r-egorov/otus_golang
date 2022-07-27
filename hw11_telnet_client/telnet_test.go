package main

import (
	"bytes"
	"io/ioutil"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type testCase struct {
	name, sendInTelnet, sendFromTelnet string
}

func TestTelnetClient(t *testing.T) {
	testCases := []testCase{
		{
			name:           "basic",
			sendFromTelnet: "hello\n",
			sendInTelnet:   "world\n",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			l, err := net.Listen("tcp", "127.0.0.1:")
			require.NoError(t, err)
			defer func() { require.NoError(t, l.Close()) }()

			var wg sync.WaitGroup
			wg.Add(2)

			go func() {
				defer wg.Done()

				in := &bytes.Buffer{}
				out := &bytes.Buffer{}

				timeout, err := time.ParseDuration("10s")
				require.NoError(t, err)

				client := NewTelnetClient(l.Addr().String(), timeout, ioutil.NopCloser(in), out)
				require.NoError(t, client.Connect())
				defer func() { require.NoError(t, client.Close()) }()

				in.WriteString(tc.sendFromTelnet)
				err = client.Send()
				require.NoError(t, err)

				err = client.Receive()
				require.NoError(t, err)
				require.Equal(t, tc.sendInTelnet, out.String())
			}()

			go func() {
				defer wg.Done()

				conn, err := l.Accept()
				require.NoError(t, err)
				require.NotNil(t, conn)
				defer func() { require.NoError(t, conn.Close()) }()

				request := make([]byte, 1024)
				n, err := conn.Read(request)
				require.NoError(t, err)
				require.Equal(t, tc.sendFromTelnet, string(request)[:n])

				n, err = conn.Write([]byte(tc.sendInTelnet))
				require.NoError(t, err)
				require.NotEqual(t, 0, n)
			}()

			wg.Wait()
		})
	}
}
