package limit

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tech-xiwi/limit/redis"
)

func TestPeriodLimit_Take(t *testing.T) {
	testPeriodLimit(t)
}

func TestPeriodLimit_TakeWithAlign(t *testing.T) {
	testPeriodLimit(t, Align())
}

func TestPeriodLimit_RedisUnavailable(t *testing.T) {
	const (
		seconds = 1
		quota   = 5
	)
	l := NewPeriodLimit(seconds, quota, redis.New("localhost:6379"), "periodlimit")
	val, err := l.Take("first")
	t.Log(val, err)
	assert.NotNil(t, err)
	assert.Equal(t, 0, val)
}

func testPeriodLimit(t *testing.T, opts ...PeriodOption) {
	store, clean, err := CreateRedis()
	assert.Nil(t, err)
	defer clean()

	const (
		seconds = 1
		total   = 100
		quota   = 5
	)
	l := NewPeriodLimit(seconds, quota, store, "periodlimit", opts...)
	var allowed, hitQuota, overQuota int
	for i := 0; i < total; i++ {
		val, err := l.Take("first")
		if err != nil {
			t.Error(err)
		}
		switch val {
		case Allowed:
			allowed++
		case HitQuota:
			hitQuota++
		case OverQuota:
			overQuota++
		default:
			t.Error("unknown status")
		}
	}

	assert.Equal(t, quota-1, allowed)
	assert.Equal(t, 1, hitQuota)
	assert.Equal(t, total-quota, overQuota)
}

// CreateRedis returns a in process redis.Redis.
func CreateRedis() (r *redis.Redis, clean func(), err error) {
	return redis.New("localhost:6379"), func() {
		ch := make(chan PlaceholderType)
		go func() {
			close(ch)
		}()

		select {
		case <-ch:
		case <-time.After(time.Second):
		}
	}, nil
}

// Placeholder is a placeholder object that can be used globally.
var Placeholder PlaceholderType

type (
	// AnyType can be used to hold any type.
	AnyType = interface{}
	// PlaceholderType represents a placeholder type.
	PlaceholderType = struct{}
)
