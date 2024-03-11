package service

import (
	"applicationDesignTest/helpers"
	"applicationDesignTest/internal/models"
	"applicationDesignTest/internal/models/dto"
	"applicationDesignTest/pkg/log"
	"context"
	"sync"
	"time"
)

type roomService interface {
	ReserveRooms(availability []models.RoomAvailability,
		daysToBook []time.Time,
		order dto.Order) error
	GetAvailabilityRooms(ctx context.Context) ([]models.RoomAvailability, error)
	UpdateAvailabilityRooms(ctx context.Context, availability []models.RoomAvailability) error
}

type orderRepo interface {
	CreateOrder(ctx context.Context, order models.Order) (id int, err error)
	DeleteOrder(ctx context.Context, id int) error
}

type OrderService struct {
	orderRepo   orderRepo
	roomService roomService
	logger      log.Logging
	sync.Mutex
}

func NewOrder(orderRepo orderRepo, roomService roomService, logger log.Logging) *OrderService {
	return &OrderService{
		orderRepo:   orderRepo,
		roomService: roomService,
		logger:      logger,
	}
}

func (o *OrderService) CreateOrder(ctx context.Context, orderRqt dto.Order) (*models.Order, error) {
	order := models.Order{HotelID: orderRqt.HotelID,
		RoomID:    orderRqt.RoomID,
		UserEmail: orderRqt.UserEmail,
		From:      orderRqt.From,
		To:        orderRqt.To,
	}

	daysToBook := helpers.DaysBetween(order.From, order.To)
	unavailableDays := make(map[time.Time]struct{}, len(daysToBook))
	for _, day := range daysToBook {
		unavailableDays[day] = struct{}{}
	}

	//TODO: не лучшее решение - используем блокировки мютексом на уровне сервиса, так как у нас in-memory хранилище с отсутствием блокировок и транзакций.
	//TODO: также отсутствует транзакция, по той же причине.
	o.Lock()
	defer o.Unlock()

	availability, err := o.roomService.GetAvailabilityRooms(ctx)
	if err != nil {
		o.logger.LogErrorf("error get availability rooms: %s", err)
		return nil, err
	}

	if err = o.roomService.ReserveRooms(availability, daysToBook, orderRqt); err != nil {
		o.logger.LogErrorf("error reservation rooms: %s", err)
		return nil, err
	}

	orderId, err := o.orderRepo.CreateOrder(ctx, order)
	if err != nil {
		o.logger.LogErrorf("error create order: %s", err)
		return nil, err
	}

	if err = o.roomService.UpdateAvailabilityRooms(ctx, availability); err != nil {
		o.logger.LogErrorf("error update availability rooms: %s", err)
		if deleteErr := o.orderRepo.DeleteOrder(ctx, orderId); deleteErr != nil {
			o.logger.LogErrorf("error rolling back order creation: %s", deleteErr)
		}
		return nil, err
	}

	order.ID = orderId
	o.logger.LogInfo("success create order for: %s", order.UserEmail)

	return &order, err
}
