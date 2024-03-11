package service

import (
	"applicationDesignTest/helpers"
	"applicationDesignTest/internal/models"
	"applicationDesignTest/internal/models/dto"
	"applicationDesignTest/pkg/log"
	"context"
	"errors"
	"sync"
)

type roomService interface {
	ReserveRoom(ctx context.Context, reserveRoom dto.ReserveRooms) error
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

func NewOrderService(orderRepo orderRepo, roomService roomService, logger log.Logging) *OrderService {
	return &OrderService{
		orderRepo:   orderRepo,
		roomService: roomService,
		logger:      logger,
	}
}

func (o *OrderService) CreateOrder(ctx context.Context, orderRqt dto.Order) (order *models.Order, err error) {
	var orderId int
	order = &models.Order{HotelID: orderRqt.HotelID,
		RoomID:    orderRqt.RoomID,
		UserEmail: orderRqt.UserEmail,
		From:      orderRqt.From,
		To:        orderRqt.To,
	}
	daysToBook := helpers.DaysBetween(order.From, order.To)

	//TODO: не лучшее решение - используем блокировки мютексом на уровне сервиса, так как у нас in-memory хранилище с отсутствием блокировок и транзакций.
	//TODO: также отсутствует транзакция, по той же причине.
	o.Lock()
	defer o.Unlock()

	//Механизм роллбэка,в случае непредвиденной ошибки
	availabilitySnapshot, err := o.roomService.GetAvailabilityRooms(ctx)
	defer func() {
		if err != nil {
			if errors.Is(err, models.ErrNotAvailableRooms) ||
				errors.Is(err, models.ErrNotAvailableRooms) {
				return
			}
			o.logger.LogInfo("start rollback availability")
			err = o.roomService.UpdateAvailabilityRooms(ctx, availabilitySnapshot)
			if err != nil {
				o.logger.LogErrorf("error rollback availability: %s", err)
				return
			}

			if orderId != 0 {
				o.logger.LogInfo("start rollback order id")
				if err = o.orderRepo.DeleteOrder(ctx, orderId); err != nil {
					o.logger.LogErrorf("error rollback availability: %s", err)
					return
				}
			}
		}
	}()

	reserveRoomRqt := dto.ReserveRooms{HotelID: orderRqt.HotelID,
		RoomID: orderRqt.RoomID,
		Dates:  daysToBook}
	if err := o.roomService.ReserveRoom(ctx, reserveRoomRqt); err != nil {
		o.logger.LogErrorf("error reservation rooms: %s", err)

		return nil, err
	}

	orderId, err = o.orderRepo.CreateOrder(ctx, *order)
	if err != nil {
		o.logger.LogErrorf("error create order: %s", err)
		return nil, err
	}

	order.ID = orderId
	o.logger.LogInfo("success create order for: %s", order.UserEmail)

	return order, err
}
