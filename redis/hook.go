package redis

import (
	"context"
	red "github.com/go-redis/redis/v8"
	"github.com/tech-xiwi/limit/utils"
	"log"
	"strings"
	"time"
)

var (
	startTimeKey = contextKey("startTime")
	durationHook = hook{}
)

type (
	contextKey string
	hook       struct{}
)

func (h hook) BeforeProcess(ctx context.Context, _ red.Cmder) (context.Context, error) {
	return context.WithValue(ctx, startTimeKey, utils.Now()), nil
}

func (h hook) AfterProcess(ctx context.Context, cmd red.Cmder) error {
	val := ctx.Value(startTimeKey)
	if val == nil {
		return nil
	}

	start, ok := val.(time.Duration)
	if !ok {
		return nil
	}

	duration := utils.Since(start)
	if duration > slowThreshold.Load() {
		logDuration(ctx, cmd, duration)
	}

	return nil
}

func (h hook) BeforeProcessPipeline(ctx context.Context, _ []red.Cmder) (context.Context, error) {
	return context.WithValue(ctx, startTimeKey, utils.Now()), nil
}

func (h hook) AfterProcessPipeline(ctx context.Context, cmds []red.Cmder) error {
	if len(cmds) == 0 {
		return nil
	}

	val := ctx.Value(startTimeKey)
	if val == nil {
		return nil
	}

	start, ok := val.(time.Duration)
	if !ok {
		return nil
	}

	duration := utils.Since(start)
	if duration > slowThreshold.Load()*time.Duration(len(cmds)) {
		logDuration(ctx, cmds[0], duration)
	}

	return nil
}

func logDuration(ctx context.Context, cmd red.Cmder, duration time.Duration) {
	var buf strings.Builder
	for i, arg := range cmd.Args() {
		if i > 0 {
			buf.WriteByte(' ')
		}
		buf.WriteString(utils.Repr(arg))
	}
	log.Printf("[REDIS] slowcall on executing: %s", buf.String())
}
