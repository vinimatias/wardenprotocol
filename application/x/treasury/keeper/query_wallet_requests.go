package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"gitlab.qredo.com/qrdochain/fusionchain/x/treasury/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) WalletRequests(goCtx context.Context, req *types.QueryWalletRequestsRequest) (*types.QueryWalletRequestsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	store := ctx.KVStore(k.storeKey)
	workspaceStore := prefix.NewStore(store, types.KeyPrefix(types.WalletRequestKey))

	walletRequests, pageRes, err := query.GenericFilteredPaginate(k.cdc, workspaceStore, req.Pagination, func(key []byte, value *types.WalletRequest) (*types.WalletRequest, error) {
		if req.XStatus == nil {
			return value, nil
		}

		reqStatus := req.XStatus.(*types.QueryWalletRequestsRequest_Status).Status
		if value.Status != reqStatus {
			return nil, nil
		}

		return value, nil
	}, func() *types.WalletRequest { return &types.WalletRequest{} })

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryWalletRequestsResponse{
		WalletRequests: walletRequests,
		Pagination:     pageRes,
	}, nil
}
