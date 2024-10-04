package pubsub

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPubSub(t *testing.T) {
	t.Parallel()

	ps := NewPubSub()

	ctx := context.Background()
	ch := ps.Subscribe("topic", ctx, 10)
	assert.NotNil(t, ch)

	ps.Publish("topic", "msg")
	msg, ok := <-ch
	assert.True(t, ok)
	assert.Equal(t, "msg", msg)

	// unsubscribe
	ps.Unsubscribe("topic", ctx)
	assert.Nil(t, ps.subs["topic"][ctx])

	_, ok = <-ch
	assert.False(t, ok, "channel should be closed")
	assert.Nil(t, ps.subs["topic"][ctx])

	ps.Publish("topic", "msg2")
}
