package redis

import (
	. "testing"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/pmoroney/go-source/event"
)

type TestEvent struct {
	Foo string
}

// Make sure we can record events
func TestRedisRecord(t *T) {
	event.Register(TestEvent{})
	event := event.EventMessage{
		Data: TestEvent{
			Foo: "bar",
		},
		ID:        event.NewID(),
		SeqID:     1,
		Timestamp: time.Now(),
		EventType: "TestEvent",
	}

	redisTestServer := redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ":6379")
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	recoder := NewRedisEventStore(redisTestServer)
	recoder.Record(event)

	events, err := recoder.GetEvents(event.ID)
	if err != nil {
		t.Fatal(err)
	}

	if len(events) != 1 {
		t.Fatal("Events Length != 1")
	}
}

// Make sure we can record events, even if they are the same
func TestRedisRecordSameEventTwice(t *T) {
	event.Register(TestEvent{})
	event := event.EventMessage{
		Data: TestEvent{
			Foo: "bar",
		},
		ID:        event.NewID(),
		SeqID:     1,
		Timestamp: time.Now(),
		EventType: "TestEvent",
	}

	redisTestServer := redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ":6379")
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	recoder := NewRedisEventStore(redisTestServer)
	recoder.Record(event)
	event.SeqID = 2
	recoder.Record(event)

	events, err := recoder.GetEvents(event.ID)
	if err != nil {
		t.Fatal(err)
	}

	if len(events) != 2 {
		t.Fatal("Events Length != 2")
	}
}
