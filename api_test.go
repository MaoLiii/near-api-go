package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/textileio/near-api-go/keys"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/require"
	"github.com/textileio/near-api-go/types"

	"testing"
)

var ctx = context.Background()

func TestIt(t *testing.T) {
	c, cleanup := makeClient(t)
	defer cleanup()
	require.NotNil(t, c)
}

// func TestViewCode(t *testing.T) {
// 	c, cleanup := makeClient(t)
// 	defer cleanup()
// 	res, err := c.ViewCode(ctx, "filecoin-bridge.testnet")
// 	require.NoError(t, err)
// 	require.NotNil(t, res)
// }

// func TestDeployContract(t *testing.T) {
// 	c, cleanup := makeClient(t)
// 	defer cleanup()
// 	res, err := c.ViewCode(ctx, "filecoin-bridge.testnet")
// 	require.NoError(t, err)
// 	require.NotNil(t, res)

// 	bytes, err := base64.StdEncoding.DecodeString(res.CodeBase64)
// 	require.NoError(t, err)
// 	require.NotEmpty(t, bytes)

// 	outcome, err := c.Account("<account id>").DeployContract(ctx, bytes)
// 	require.NoError(t, err)
// 	require.NotNil(t, outcome)

// 	res2, err := c.ViewCode(ctx, "<account id>")
// 	require.NoError(t, err)
// 	require.NotNil(t, res)

// 	require.Equal(t, res.Hash, res2.Hash)
// }

// func TestDataChanges(t *testing.T) {
// 	c, cleanup := makeClient(t)
// 	defer cleanup()
// 	res, err := c.DataChanges(ctx, []string{"filecoin-bridge.testnet"}, DataChangesWithFinality("final"))
// 	require.NoError(t, err)
// 	require.NotNil(t, res)
// }

func makeClient(t *testing.T) (*Client, func()) {
	rpcClient, err := rpc.DialContext(ctx, "https://rpc.testnet.near.org")
	require.NoError(t, err)

	// keys, err := keys.NewKeyPairFromString(
	// 	"ed25519:xxxx",
	// )
	// require.NoError(t, err)

	config := &types.Config{
		RPCClient: rpcClient,
		// Signer:    keys,
		NetworkID: "testnet",
	}
	c, err := NewClient(config)

	require.NoError(t, err)
	return c, func() {
		rpcClient.Close()
	}
}

func TestGetTokenInfo(t *testing.T) {
	var (
		nodeUrl          = "https://rpc.mainnet.near.org"
		walletPrivateKey = "ed25519:4v9CciSAsKxRpyKqtr28fGiEHLCFH2F7y5do8M968hKNbN84fG9nAeCyuHqpEWYk1UBQyGXJfoQj9vKUgdqVx4zg" // 钱包私钥
		networkId        = "mainnet"
		nearRpc          *rpc.Client
		nearApi          *Client
		err              error
	)
	nearRpc, nearApi, err = InitNearClient(nodeUrl, walletPrivateKey, networkId)
	if err != nil {
		fmt.Println(err)
		return
	}
	_ = nearRpc

	r, e := GetTokenInfo(nearApi, "game.hot.tg")
	if e != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(r)
}

func InitNearClient(nodeUrl string, walletPrivateKey string, networkId string) (nearRpc *rpc.Client, nearClient *Client, err error) {
	var (
		rpcClient *rpc.Client
		nearKeys  keys.KeyPair
	)
	if len(nodeUrl) == 0 || len(walletPrivateKey) == 0 || len(networkId) == 0 {
		return
	}

	customHttpClient := &http.Client{
		Timeout: time.Second * 60, // 例如，30秒超时
	}
	rpcClient, err = rpc.DialHTTPWithClient(nodeUrl, customHttpClient)
	if err != nil {
		return
	}
	rpcClient.SetHeader("Accept", "*/*")
	rpcClient.SetHeader("Referer", "https://app.ref.finance/")
	nearRpc = rpcClient

	nearKeys, err = keys.NewKeyPairFromString(walletPrivateKey)
	if err != nil {
		return
	}
	nearConfig := &types.Config{
		RPCClient: rpcClient,
		NetworkID: networkId,
		Signer:    nearKeys,
	}
	nearClient, err = NewClient(nearConfig)
	if err != nil {
		return
	}
	return
}

func GetTokenInfo(nearClient *Client, tokenCode string) (resMap map[string]interface{}, err error) {
	var (
		res *CallFunctionResponse
	)
	res, err = nearClient.CallFunction(ctx, tokenCode, "ft_metadata", CallFunctionWithFinality("final"))
	if err != nil {
		return
	}
	err = json.Unmarshal(res.Result, &resMap)
	if err != nil {
		return
	}
	return
}
