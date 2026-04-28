package worker

import (
	"citra-wastra-be/dto"
	"citra-wastra-be/repository"
	"citra-wastra-be/service"
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

const maxRetry = 3

type XPWorker struct {
	service   service.GamificationService
	queueRepo repository.QueueRepository
}

func NewXPWorker(s service.GamificationService, q repository.QueueRepository) *XPWorker {
	return &XPWorker{
		service:   s,
		queueRepo: q,
	}
}

func (w *XPWorker) Start(ctx context.Context) {
	log.Println("xp worker is running and waiting for jobs..")

	for {
		select {
		case <-ctx.Done():
			log.Println("xp worker stopped")
			return
		default:
		}
		rawData, err := w.queueRepo.DequeueXPJob(ctx)
		if err != nil {
			if errors.Is(err, redis.Nil) {
				time.Sleep(2 * time.Second)
				continue
			}
			log.Printf("worker error: %v", err)
			time.Sleep(1 * time.Second)
			continue
		}

		var payload dto.XPJobPayload

		if err := json.Unmarshal([]byte(rawData), &payload); err != nil {
			log.Printf("invalid payload: %v", err)

			_ = w.queueRepo.AckJob(ctx, rawData)
			continue
		}

		log.Printf("processing XP: User %s (+%d) retry=%d",
			payload.UserID, payload.XPGain, payload.Retry)

		err = w.service.ProcessXPAndBadges(payload)
		if err != nil {
			log.Printf("failed: %v", err)

			payload.Retry++

			if payload.Retry < maxRetry {
				log.Printf("retrying job user=%s retry=%d", payload.UserID, payload.Retry)

				newData, _ := json.Marshal(payload)
				_ = w.queueRepo.EnqueueRaw(ctx, newData)
			} else {
				log.Printf("job moved to dead letter (max retry reached)")
				_ = w.queueRepo.AckJob(ctx, rawData)
			}
			continue
		}

		if err := w.queueRepo.AckJob(ctx, rawData); err != nil {
			log.Printf("failed to ACK job: %v", err)
		} else {
			log.Printf("job success user=%s", payload.UserID)
		}
	}
}
