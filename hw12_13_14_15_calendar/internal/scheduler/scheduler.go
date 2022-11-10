package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/novikov-ai/practice-misc/hw12_13_14_15_calendar/configs"
	"github.com/novikov-ai/practice-misc/hw12_13_14_15_calendar/internal/app"
	"github.com/novikov-ai/practice-misc/hw12_13_14_15_calendar/pkg/ampq"
	"github.com/novikov-ai/practice-misc/hw12_13_14_15_calendar/pkg/models"
)

type Scheduler struct {
	storage         app.Storage
	logger          app.Logger
	config          configs.ConfigScheduler
	messageBroker   ampq.Client
	processedEvents map[string]struct{} // for excluding duplicates
	scanEvery       time.Duration
}

const (
	QuerySelectEventIDsToNotify = `SELECT * FROM events
WHERE NOW() = (date - notified_before)`

	QueryDeleteOldEvents = `DELETE FROM events
WHERE EXTRACT(YEAR FROM NOW()) - (EXTRACT(YEAR FROM date)) > 1`
)

func New(st app.Storage, log app.Logger, conf configs.ConfigScheduler, broker ampq.Client) *Scheduler {
	newScheduler := Scheduler{
		storage: st, logger: log, config: conf,
		messageBroker:   broker,
		processedEvents: map[string]struct{}{},
		scanEvery:       time.Second * 5,
	}

	frequency, err := time.ParseDuration(conf.Scheduler.ScanFrequency)
	if err == nil {
		newScheduler.scanEvery = frequency
	}

	return &newScheduler
}

func (sch *Scheduler) Scan(ctx context.Context) error {
	wg := sync.WaitGroup{}
	errors := make(chan error, 1)

	for {
		select {
		case <-ctx.Done():
			return nil

		case err := <-errors:
			sch.logger.Error(schedulerLog(err.Error()))
			return err

		case <-time.After(sch.scanEvery):
			wg.Add(1)
			go func() {
				defer wg.Done()

				err := sendNotifications(ctx, sch, QuerySelectEventIDsToNotify)
				if err != nil {
					errors <- err
				}
			}()

			wg.Add(1)
			go func() {
				defer wg.Done()

				err := deleteOldEvents(ctx, sch)
				if err != nil {
					errors <- err
				}
			}()

			wg.Wait()
		}
	}
}

func sendNotifications(ctx context.Context, sch *Scheduler, query string) error {
	rows, err := sch.storage.RunQuery(ctx, query)
	if err != nil {
		return err
	}

	notifications := make([]models.Notification, 0)
	for rows.Next() {
		var (
			id, title, description, userID, duration, notifiedBefore string
			date                                                     time.Time
		)

		err := rows.Scan(&id, &title, &description, &userID, &date, &duration, &notifiedBefore)
		if err != nil {
			return err
		}

		_, ok := sch.processedEvents[id]
		if ok {
			continue
		}

		notification := models.Notification{
			EventID: id,
			Title:   title,
			UserID:  userID,
			Date:    date,
		}
		notifications = append(notifications, notification)
	}
	defer rows.Close()

	for _, n := range notifications {
		encoded, err := json.Marshal(n)
		if err != nil {
			return err
		}

		err = sch.messageBroker.Send(ctx, string(encoded))
		if err != nil {
			sch.logger.Warn(fmt.Sprintf("Notification for event: %s wasn't sent. Error: %s\n", n.EventID, err.Error()))
			return err
		}

		sch.processedEvents[n.EventID] = struct{}{}
	}

	return nil
}

func deleteOldEvents(ctx context.Context, sch *Scheduler) error {
	rows, err := sch.storage.RunQuery(ctx, QueryDeleteOldEvents)
	if err != nil {
		return err
	}
	defer rows.Close()

	sch.logger.Info("cleaned up old events")

	return nil
}

func schedulerLog(message string) string {
	return fmt.Sprintf("Scheduler: %s\n", message)
}
