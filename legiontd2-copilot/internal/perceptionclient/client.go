package perceptionclient

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	pb "github.com/yourname/legiontd2-copilot/internal/perceptionclient/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type EconomyState struct {
	Mythium         int
	Income          int
	WaveNumber      int
	WaveTimerSec    int
	KingHPPercent   int
	AllyKingHP      int
	Confidence      float32
}

type Client struct {
	conn   *grpc.ClientConn
	client pb.PerceptionServiceClient
}

func New(address string) (*Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("grpc dial: %w", err)
	}

	slog.Info("perception client connected", "address", address)
	return &Client{
		conn:   conn,
		client: pb.NewPerceptionServiceClient(conn),
	}, nil
}

func (c *Client) ReadEconomy(ctx context.Context) (*EconomyState, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	resp, err := c.client.ReadEconomy(ctx, &pb.ReadEconomyRequest{})
	if err != nil {
		return nil, fmt.Errorf("read economy: %w", err)
	}

	return &EconomyState{
		Mythium:       int(resp.Mythium),
		Income:        int(resp.Income),
		WaveNumber:    int(resp.WaveNumber),
		WaveTimerSec:  int(resp.WaveTimerSeconds),
		KingHPPercent: int(resp.KingHpPercent),
		AllyKingHP:    int(resp.AllyKingHpPercent),
		Confidence:    resp.Confidence,
	}, nil
}

func (c *Client) HealthCheck(ctx context.Context) (bool, bool, string) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	resp, err := c.client.HealthCheck(ctx, &pb.HealthCheckRequest{})
	if err != nil {
		return false, false, err.Error()
	}

	return resp.GameWindowDetected, resp.CaptureHealthy, resp.Message
}

func (c *Client) Close() error {
	return c.conn.Close()
}
