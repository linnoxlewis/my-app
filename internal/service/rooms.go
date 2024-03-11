package service

import (
	"applicationDesignTest/internal/models"
	"applicationDesignTest/internal/models/dto"
	"applicationDesignTest/pkg/log"
	"context"
	"fmt"
	"time"
)

type roomRepo interface {
	GetAvailabilityRooms(ctx context.Context) ([]models.RoomAvailability, error)
	UpdateAvailabilityRooms(ctx context.Context, roomAvailability []models.RoomAvailability) error
}

type RoomService struct {
	roomRepo roomRepo
	logger   log.Logging
}

func NewRoomService(roomRepo roomRepo, logger log.Logging) *RoomService {
	return &RoomService{roomRepo: roomRepo,
		logger: logger,
	}
}

func (r *RoomService) ReserveRooms(availability []models.RoomAvailability,
	daysToBook []time.Time,
	order dto.Order) error {
	availabilityMap := make(map[string]int, len(availability))

	getHashKey := func(hotelId, roomId string, date time.Time) string {
		return fmt.Sprintf("%s-%s-%s", hotelId, roomId, date)
	}

	for key, value := range availability {
		hashKey := getHashKey(value.HotelID, value.RoomID, value.Date)
		availabilityMap[hashKey] = key
	}

	for _, dayToBook := range daysToBook {
		key := getHashKey(order.HotelID, order.RoomID, dayToBook)
		if idx, exists := availabilityMap[key]; exists {
			if availability[idx].Quota > 0 {
				availability[idx].Quota--
			} else {
				return models.ErrNotAvailableRooms
			}
		} else {
			return models.ErrNotFoundInformation
		}
	}

	return nil
}

func (r *RoomService) GetAvailabilityRooms(ctx context.Context) ([]models.RoomAvailability, error) {
	availability, err := r.roomRepo.GetAvailabilityRooms(ctx)
	if err != nil {
		r.logger.LogErrorf("error get availability rooms: %s", err)
		return nil, err
	}

	return availability, nil
}

func (r *RoomService) UpdateAvailabilityRooms(ctx context.Context, availability []models.RoomAvailability) error {
	if err := r.roomRepo.UpdateAvailabilityRooms(ctx, availability); err != nil {
		r.logger.LogErrorf("error get availability rooms: %s", err)
		return err
	}

	return nil
}
