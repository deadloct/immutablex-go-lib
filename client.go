package lib

import (
	"os"
	"sync"

	"github.com/immutable/imx-core-sdk-golang/imx"
	"github.com/immutable/imx-core-sdk-golang/imx/api"
	log "github.com/sirupsen/logrus"
)

const MaxAssetsPerReq = 200

var AlchemyKey string

func init() {
	AlchemyKey = os.Getenv("ALCHEMY_API_KEY")
	if AlchemyKey == "" {
		log.Panic("no alchemy api key provided, get one at alchemy.com")
	}
}

type IMXClientWrapper interface {
	Start() error
	Stop()
	GetClient() *imx.Client
}

type Client struct {
	key       string
	imxClient *imx.Client

	sync.Mutex
}

func NewClient() *Client {
	return &Client{key: AlchemyKey}
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
