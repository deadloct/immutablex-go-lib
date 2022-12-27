package lib

import (
	"sync"

	"github.com/immutable/imx-core-sdk-golang/imx"
	"github.com/immutable/imx-core-sdk-golang/imx/api"
)

type ClientWrapper interface {
	Start() error
	Stop()
	GetClient() *imx.Client
}

type Client struct {
	key       string
	imxClient *imx.Client

	sync.Mutex
}

func NewClient(alchemyKey string) *Client {
	return &Client{key: alchemyKey}
}

func (c *Client) Start() error {
	if c.imxClient != nil {
		return nil
	}

	cfg := imx.Config{
		AlchemyAPIKey: c.key,
		APIConfig:     api.NewConfiguration(),
		Environment:   imx.Mainnet,
	}

	client, err := imx.NewClient(&cfg)
	if err != nil {
		return err
	}

	if client != nil {
		c.Lock()
		c.imxClient = client
		c.Unlock()
	}

	return nil
}

func (c *Client) Stop() {
	c.Lock()
	defer c.Unlock()

	if c.imxClient == nil {
		return
	}

	c.imxClient.EthClient.Close()
	c.imxClient = nil
}

func (c *Client) GetClient() *imx.Client {
	c.Lock()
	defer c.Unlock()

	return c.imxClient
}
