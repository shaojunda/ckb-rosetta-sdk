package services

import (
	"context"
	"github.com/nervosnetwork/ckb-rosetta-sdk/server/config"

	"github.com/coinbase/rosetta-sdk-go/server"
	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/nervosnetwork/ckb-sdk-go/rpc"
)

// NetworkAPIService implements the server.NetworkAPIService interface.
type NetworkAPIService struct {
	network *types.NetworkIdentifier
	client  rpc.Client
	cfg     *config.Config
}

// NewNetworkAPIService creates a new instance of a NetworkAPIService.
func NewNetworkAPIService(network *types.NetworkIdentifier, client rpc.Client, cfg *config.Config) server.NetworkAPIServicer {
	return &NetworkAPIService{
		network: network,
		client:  client,
		cfg:     cfg,
	}
}

// NetworkList implements the /network/list endpoint
func (s *NetworkAPIService) NetworkList(
	ctx context.Context,
	request *types.MetadataRequest,
) (*types.NetworkListResponse, *types.Error) {
	return &types.NetworkListResponse{
		NetworkIdentifiers: []*types.NetworkIdentifier{
			s.network,
		},
	}, nil
}

// NetworkStatus implements the /network/status endpoint.
func (s *NetworkAPIService) NetworkStatus(
	ctx context.Context,
	request *types.NetworkRequest,
) (*types.NetworkStatusResponse, *types.Error) {
	genesis, err := s.client.GetHeaderByNumber(context.Background(), 0)
	if err != nil {
		return nil, RpcError
	}
	peers, err := s.client.GetPeers(context.Background())
	if err != nil {
		return nil, RpcError
	}
	header, err := s.client.GetTip(context.Background())
	if err != nil {
		return nil, RpcError
	}
	nodeHeader, err := s.client.GetHeaderByNumber(context.Background(), header.BlockNumber)
	if err != nil {
		return nil, RpcError
	}

	result := &types.NetworkStatusResponse{
		CurrentBlockIdentifier: &types.BlockIdentifier{
			Index: int64(nodeHeader.Number),
			Hash:  nodeHeader.Hash.String(),
		},
		CurrentBlockTimestamp: int64(nodeHeader.Timestamp),
		GenesisBlockIdentifier: &types.BlockIdentifier{
			Index: 0,
			Hash:  genesis.Hash.String(),
		},
		Peers: []*types.Peer{},
	}

	for _, peer := range peers {
		result.Peers = append(result.Peers, &types.Peer{
			PeerID: peer.NodeId,
		})
	}

	return result, nil
}

// NetworkOptions implements the /network/options endpoint.
func (s *NetworkAPIService) NetworkOptions(
	ctx context.Context,
	request *types.NetworkRequest,
) (*types.NetworkOptionsResponse, *types.Error) {
	node, err := s.client.LocalNodeInfo(context.Background())
	if err != nil {
		return nil, RpcError
	}

	return &types.NetworkOptionsResponse{
		Version: &types.Version{
			RosettaVersion: "1.4.5",
			NodeVersion:    node.Version,
		},
		Allow: &types.Allow{
			OperationStatuses: []*types.OperationStatus{
				{
					Status:     "Success",
					Successful: true,
				},
			},
			OperationTypes: SupportedOperationTypes,
			Errors:         AllErrorTypes,
		},
	}, nil
}
