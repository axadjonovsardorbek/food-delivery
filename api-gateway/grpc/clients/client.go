package clients

import (
	"gateway/config"
	cp "gateway/genproto/order"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcClients struct {
	Product          cp.ProductServiceClient
	Cart           cp.CartServiceClient
	// SharedMemory    cp.SharedMemoriesServiceClient
	// Comment         cp.CommentsServiceClient
	// Milestone       cp.MilestonesServiceClient
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

	// connT, err := grpc.NewClient(cfg.TIMELINE_HOST+cfg.TIMELINE_PORT,
	// 	grpc.WithTransportCredentials(insecure.NewCredentials()))
	// if err != nil {
	// 	return nil, err
	// }

	return &GrpcClients{
		Product:          cp.NewProductServiceClient(connO),
		Cart:           cp.NewCartServiceClient(connO),
		// SharedMemory:    cp.NewSharedMemoriesServiceClient(connM),
		// Comment:         cp.NewCommentsServiceClient(connM),
		// Milestone:       cp.NewMilestonesServiceClient(connT),
		// CustomEvent:     cp.NewCustomEventsServiceClient(connT),
		// PersonalEvent:   cp.NewPersonalEventsServiceClient(connT),
		// HistoricalEvent: cp.NewHistoricalEventsServiceClient(connT),
	}, nil
}
