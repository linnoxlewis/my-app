package repository

import (
	"applicationDesignTest/helpers"
	"applicationDesignTest/internal/models"
	"context"
	"sync"
)

var availability = []models.RoomAvailability{
	{"reddison", "lux", helpers.Date(2024, 1, 1), 1},
	{"reddison", "lux", helpers.Date(2024, 1, 2), 1},
	{"reddison", "lux", helpers.Date(2024, 1, 3), 1},
	{"reddison", "lux", helpers.Date(2024, 1, 4), 1},
	{"reddison", "lux", helpers.Date(2024, 1, 5), 0},
}

type RoomRepo struct {
	availability []models.RoomAvailability
	sync.Mutex
}

func NewRoom() *RoomRepo {
	return &RoomRepo{
		availability: availability,
	}
}

func (r *RoomRepo) GetAvailabilityRooms(ctx context.Context) ([]models.RoomAvailability, error) {
	result := make([]models.RoomAvailability, len(r.availability), len(r.availability))
	copy(result, r.availability)

	return result, nil
}

func (r *RoomRepo) UpdateAvailabilityRooms(ctx context.Context, roomAvailability []models.RoomAvailability) error {
	r.Lock()
	defer r.Unlock()
	r.availability = roomAvailability

	return nil
}
