package internal

// ERC20Client is used to be able to get the Ethereum Blockchain syncing status,
type ERC20Client struct {
	client       *CommandLineClient
	rootEndpoint string
}

// NewERC20Client creates a new ERC20Client,
// that can be used for easy interaction with the ERC20 API exposed via the ERC20 endpoints
func NewERC20Client(cli *CommandLineClient) *ERC20Client {
	if cli == nil {
		panic("no CommandLineClient given")
	}
	return &ERC20Client{
		client:       cli,
		rootEndpoint: "/erc20",
	}
}
