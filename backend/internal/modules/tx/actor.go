package tx

import "github.com/5nat/nft-auction-platform/backend/internal/modules/auth"

// Actor 表示当前请求交易构造的用户身份。
// TxService 不直接依赖 gin.Context，也不直接解析 JWT。
// Handler 负责把 AuthMiddleware 中的 CurrentUser 转换成 Actor 传进来。
type Actor struct {
	UserID        uint64
	WalletID      uint64
	WalletAddress string
	ChainID       int64
}

func ActorFromAuth(currentUser *auth.CurrentUser) Actor {
	if currentUser == nil {
		return Actor{}
	}

	return Actor{
		UserID:        currentUser.UserID,
		WalletID:      currentUser.WalletID,
		WalletAddress: currentUser.WalletAddress,
		ChainID:       currentUser.ChainID,
	}
}

func (a Actor) IsZero() bool {
	return a.UserID == 0 || a.WalletID == 0 || a.WalletAddress == "" || a.ChainID == 0
}
