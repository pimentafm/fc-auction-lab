package auction_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/pimentafm/fc-auction-lab/internal/infra/database/auction"
	"github.com/pimentafm/fc-auction-lab/internal/internal_error"
	"github.com/stretchr/testify/mock"
)

type AuctionRepositoryMock struct {
	mock.Mock
}

func (m *AuctionRepositoryMock) CloseAuctionUpdate(ctx context.Context, auctionID string) *internal_error.InternalError {
	args := m.Called(ctx, auctionID)
	return args.Get(0).(*internal_error.InternalError)
}

func TestMonitorAuction(t *testing.T) {
	repository := &AuctionRepositoryMock{}
	repository.On("CloseAuctionUpdate", context.Background(), "123").Return((*internal_error.InternalError)(nil))

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		auction.MonitorAuction(context.Background(), "123", 100*time.Millisecond, repository)
	}()
	repository.AssertNumberOfCalls(t, "CloseAuctionUpdate", 0)
	wg.Wait()
	repository.AssertNumberOfCalls(t, "CloseAuctionUpdate", 1)
}
