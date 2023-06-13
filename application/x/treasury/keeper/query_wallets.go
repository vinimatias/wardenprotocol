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

func (k Keeper) Wallets(goCtx context.Context, req *types.QueryWalletsRequest) (*types.QueryWalletsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	store := ctx.KVStore(k.storeKey)
	workspaceStore := prefix.NewStore(store, types.KeyPrefix(types.WalletKey))

	wallets, pageRes, err := query.GenericFilteredPaginate(k.cdc, workspaceStore, req.Pagination, func(key []byte, value *types.Wallet) (*types.Wallet, error) {
		if req.XWorkspaceId == nil {
			return value, nil
		}

		wID := req.XWorkspaceId.(*types.QueryWalletsRequest_WorkspaceId).WorkspaceId
		if value.WorkspaceId != wID {
			return nil, nil
		}

		return value, nil
	}, func() *types.Wallet { return &types.Wallet{} })

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryWalletsResponse{
		Wallets:    wallets,
		Pagination: pageRes,
	}, nil
}
