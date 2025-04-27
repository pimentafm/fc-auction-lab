package auction

import (
	"context"
	"os"
	"time"

	"github.com/pimentafm/fc-auction-lab/configuration/logger"
	"github.com/pimentafm/fc-auction-lab/internal/entity/auction_entity"
	"github.com/pimentafm/fc-auction-lab/internal/internal_error"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/mongo"
)

type AuctionEntityMongo struct {
	Id          string                          `bson:"_id"`
	ProductName string                          `bson:"product_name"`
	Category    string                          `bson:"category"`
	Description string                          `bson:"description"`
	Condition   auction_entity.ProductCondition `bson:"condition"`
	Status      auction_entity.AuctionStatus    `bson:"status"`
	Timestamp   int64                           `bson:"timestamp"`
}
type AuctionRepository struct {
	Collection *mongo.Collection
}

func NewAuctionRepository(database *mongo.Database) *AuctionRepository {
	return &AuctionRepository{
		Collection: database.Collection("auctions"),
	}
}

func (ar *AuctionRepository) CreateAuction(
	ctx context.Context,
	auctionEntity *auction_entity.Auction) *internal_error.InternalError {
	auctionEntityMongo := &AuctionEntityMongo{
		Id:          auctionEntity.Id,
		ProductName: auctionEntity.ProductName,
		Category:    auctionEntity.Category,
		Description: auctionEntity.Description,
		Condition:   auctionEntity.Condition,
		Status:      auctionEntity.Status,
		Timestamp:   auctionEntity.Timestamp.Unix(),
	}
	_, err := ar.Collection.InsertOne(ctx, auctionEntityMongo)
	if err != nil {
		logger.Error("Error trying to insert auction", err)
		return internal_error.NewInternalServerError("Error trying to insert auction")
	}

	go func() {
		auctionDuration := getTimeAuctionDuration()
		errMonitor := MonitorAuction(ctx, auctionEntity.Id, auctionDuration, ar)
		if errMonitor != nil {
			logger.Error("Error trying to monitor auction", errMonitor)
		}
	}()

	return nil
}

func (ar *AuctionRepository) CloseAuctionUpdate(ctx context.Context, auctionID string) *internal_error.InternalError {
	_, err := ar.Collection.UpdateOne(ctx, bson.M{"_id": auctionID}, bson.M{"$set": bson.M{"status": auction_entity.Completed}})
	if err != nil {
		logger.Error("Error trying to update auction status", err)
		return internal_error.NewInternalServerError(err.Error())
	}

	logger.Info("Auction closed", zap.String("auction_id", auctionID))

	return nil
}

type Repository interface {
	CloseAuctionUpdate(ctx context.Context, auctionID string) *internal_error.InternalError
}

func MonitorAuction(ctx context.Context, auctionID string, auctionDuration time.Duration, repository Repository) *internal_error.InternalError {
	timer := time.NewTimer(auctionDuration)

	select {
	case <-timer.C:
		err := repository.CloseAuctionUpdate(ctx, auctionID)
		if err != nil {
			return err
		}
	case <-ctx.Done():
		timer.Stop()
		return internal_error.NewInternalServerError("Auction monitoring cancelled")
	}

	return nil
}

func getTimeAuctionDuration() time.Duration {
	auctionDuration := os.Getenv("AUCTION_DURATION")
	duration, err := time.ParseDuration(auctionDuration)
	if err != nil {
		return time.Minute * 1
	}

	return duration
}
