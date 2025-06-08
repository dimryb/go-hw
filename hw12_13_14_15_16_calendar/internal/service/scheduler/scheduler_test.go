package scheduler_test

import (
	"context"
	"testing"
	"time"

	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/config"
	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/service/scheduler"
	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/types"
	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/mocks"
	"github.com/golang/mock/gomock" //nolint:depguard
	"github.com/stretchr/testify/require"
)

func TestScheduler_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockApp := mocks.NewMockApplication(ctrl)
	mockRmq := mocks.NewMockRmqClient(ctrl)
	mockLog := mocks.NewMockLogger(ctrl)

	now := time.Now()
	event := types.Event{
		ID:           "event_id",
		Title:        "Team Meeting",
		Description:  "Discuss roadmap",
		UserID:       "user1",
		StartTime:    now.Add(10 * time.Second),
		EndTime:      now.Add(1 * time.Hour),
		NotifyBefore: 600,
	}

	cfg := &config.SchedulerConfig{
		Scheduler: config.Scheduler{
			Interval:        10 * time.Millisecond,
			RetentionPeriod: 8760 * time.Hour, // 1 год
		},
	}

	mockApp.EXPECT().
		ListEventsDueBefore(gomock.Any(), gomock.Any()).
		Return([]types.Event{event}, nil).
		AnyTimes()

	mockRmq.EXPECT().
		Publish(event.UserID, gomock.Any()).
		Return(nil).
		AnyTimes()

	mockApp.EXPECT().
		DeleteOlderThan(gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()

	mockLog.EXPECT().
		Infof(gomock.Any(), gomock.Any()).
		AnyTimes()

	sched := scheduler.NewScheduler(mockApp, mockRmq, mockLog, cfg)

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	go func() {
		err := sched.Run(ctx)
		require.Error(t, err)
	}()

	time.Sleep(20 * time.Millisecond)
	cancel()
}
