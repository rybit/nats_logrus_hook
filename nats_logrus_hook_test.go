package nhook

import (
	"encoding/json"
	"log"
	"os"
	"testing"
	"time"

	"github.com/nats-io/nats"
	"github.com/nats-io/nats/test"
	"github.com/sirupsen/logrus"
	ltest "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

var nc *nats.Conn

func TestMain(m *testing.M) {
	s := test.RunDefaultServer()
	defer s.Shutdown()

	var err error
	nc, err = nats.Connect("nats://" + s.Addr().String())
	if err != nil {
		log.Fatal("Failed to connect to server: " + err.Error())
	}
	defer nc.Close()

	os.Exit(m.Run())
}

func TestSimpleSend(t *testing.T) {
	logger, logHook := ltest.NewNullLogger()

	sub, err := nc.SubscribeSync("test")
	if !assert.NoError(t, err) {
		return
	}
	defer sub.Unsubscribe()

	hook := NewNatsHook(nc, "test")

	hook.AddField("hook-level", "pikachu")
	hook.AddDynamicField("dynamic-level", func() interface{} {
		return 12
	})

	addToAll(logger, hook)

	logger.WithField("instance-level", "charizard").Info("this is a test")

	assert.Len(t, logHook.Entries, 1)
	entry := logHook.Entries[0]

	msg, err := sub.NextMsg(time.Minute)
	parsed := make(map[string]interface{})

	err = json.Unmarshal(msg.Data, &parsed)
	if !assert.NoError(t, err) {
		return
	}

	// {"dynamic-level":12,"hook-level":"pikachu","instance-level":"charizard","level":"info","msg":"this is a test","time":"2016-10-19T13:01:14-07:00"}
	assert.Len(t, parsed, 6)

	for _, m := range []map[string]interface{}{entry.Data, parsed} {
		assert.Equal(t, "pikachu", m["hook-level"])
		assert.Equal(t, "charizard", m["instance-level"])
		assert.EqualValues(t, 12, m["dynamic-level"])
	}

	assert.Equal(t, "this is a test", parsed["msg"])
	assert.Equal(t, "this is a test", entry.Message)

	assert.Equal(t, "info", parsed["level"])
	assert.Equal(t, "info", entry.Level.String())
	assert.NotEmpty(t, parsed["time"])
	assert.NotEmpty(t, entry.Time)
}

func addToAll(logger *logrus.Logger, hook logrus.Hook) {
	for l, hooks := range logger.Hooks {
		logger.Hooks[l] = append(hooks, hook)
	}
}
