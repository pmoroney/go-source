package redis

import (
	"log"

	"github.com/garyburd/redigo/redis"
	"github.com/pmoroney/go-source/event"
)

// This is a simple event recorder for testing, it doesnt actually persist the events
type RedisEventStore struct {
	pool redis.Pool
}

func NewRedisEventStore(pool redis.Pool) RedisEventStore {
	return RedisEventStore{
		pool: pool,
	}
}

func (r RedisEventStore) Record(e event.Message) {
	conn := r.pool.Get()
	defer conn.Close()

	s, err := e.Serialize()
	if err != nil {
		log.Printf("ERROR: %s\n", err)
	}

	n, err := conn.Do("RPUSH", e.ID.String(), s)
	if err != nil {
		log.Printf("ERROR: %s\n", err)
	}

	log.Printf("Recorded Event %d: %s\n", n, s)
}

func (r RedisEventStore) GetEvents(id event.ID) ([]event.Event, error) {
	conn := r.pool.Get()
	defer conn.Close()

	strings, err := redis.Strings(conn.Do("LRANGE", id.String(), 0, -1))
	if err != nil {
		return nil, err
	}

	events := make([]event.Event, len(strings))
	for i := range strings {
		log.Printf("Received Event %d: %s\n", i, strings[i])
		event, err := event.Unserialize([]byte(strings[i]))
		if err != nil {
			return nil, err
		}
		log.Printf("Parsed Event %d: %#v\n", i, event)
		events[i] = event.Data
	}

	return events, nil
}
