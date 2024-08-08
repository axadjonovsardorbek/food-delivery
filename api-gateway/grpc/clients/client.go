package clients

import (
	"gateway/config"
	cp "gateway/genproto/courier"
	op "gateway/genproto/order"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcClients struct {
	Product      op.ProductServiceClient
	Cart         op.CartServiceClient
	CartItem     op.CartItemServiceClient
	Task         cp.TaskServiceClient
	Notification cp.NotificationServiceClient
	// CustomEvent     cp.CustomEventsServiceClient
	// PersonalEvent   cp.PersonalEventsServiceClient
	// HistoricalEvent cp.HistoricalEventsServiceClient
}

func NewGrpcClients(cfg *config.Config) (*GrpcClients, error) {
	connO, err := grpc.NewClient(cfg.ORDER_HOST+cfg.ORDER_PORT,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	connC, err := grpc.NewClient(cfg.COURIER_HOST+cfg.COURIER_PORT,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &GrpcClients{
		Product:      op.NewProductServiceClient(connO),
		Cart:         op.NewCartServiceClient(connO),
		CartItem:     op.NewCartItemServiceClient(connO),
		Task:         cp.NewTaskServiceClient(connC),
		Notification: cp.NewNotificationServiceClient(connC),
		// CustomEvent:     cp.NewCustomEventsServiceClient(connT),
		// PersonalEvent:   cp.NewPersonalEventsServiceClient(connT),
		// HistoricalEvent: cp.NewHistoricalEventsServiceClient(connT),
	}, nil
}
