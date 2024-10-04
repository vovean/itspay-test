package garantexrateprovider

import (
	"context"
	"fmt"
	"io"
	"itspay/internal/entity"
	"net/http"
	"time"

	"github.com/cockroachdb/apd/v3"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type RateProvider struct {
	client http.Client
}

const requestTimeout = time.Second // TODO setup from config

func New() *RateProvider {
	return &RateProvider{
		client: http.Client{
			Transport: otelhttp.NewTransport(http.DefaultTransport),
			Timeout:   requestTimeout,
		},
	}
}

type BidsAsks struct {
	Timestamp int64 `json:"timestamp"`
	Asks      []struct {
		Price  string `json:"price"`
		Volume string `json:"volume"`
		Amount string `json:"amount"`
		Factor string `json:"factor"`
		Type   string `json:"type"`
	} `json:"asks"`
	Bids []struct {
		Price  string `json:"price"`
		Volume string `json:"volume"`
		Amount string `json:"amount"`
		Factor string `json:"factor"`
		Type   string `json:"type"`
	} `json:"bids"`
}

func (c *RateProvider) GetRate(ctx context.Context) (*entity.Rate, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://garantex.org/api/v2/depth?market=usdtrub", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create a new request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make a request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}

		return nil, fmt.Errorf("%s: %s", resp.Status, body)
	}

	bidsAsks := &BidsAsks{}
	if err := jsoniter.NewDecoder(resp.Body).Decode(bidsAsks); err != nil {
		return nil, fmt.Errorf("cannot decode response body: %w", err)
	}

	if len(bidsAsks.Asks) < 1 {
		return nil, errors.New("no asks found")
	}

	if len(bidsAsks.Bids) < 1 {
		return nil, errors.New("no bids found")
	}

	ask, _, err := apd.NewFromString(bidsAsks.Asks[0].Price)
	if err != nil {
		return nil, fmt.Errorf("cannot parse ask price (%q): %w", bidsAsks.Asks[0].Price, err)
	}

	bid, _, err := apd.NewFromString(bidsAsks.Bids[0].Price)
	if err != nil {
		return nil, fmt.Errorf("cannot parse bid price (%q): %w", bidsAsks.Bids[0].Price, err)
	}

	receivedAt := time.Unix(bidsAsks.Timestamp, 0)

	return &entity.Rate{
		Ask:        ask,
		Bid:        bid,
		ReceivedAt: receivedAt,
	}, nil
}
