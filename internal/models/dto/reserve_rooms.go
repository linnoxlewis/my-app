package dto

import "time"

type ReserveRooms struct {
	HotelID string      `json:"hotel_id"`
	RoomID  string      `json:"room_id"`
	Dates   []time.Time `json:"dates"`
}
