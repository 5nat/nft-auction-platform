package auth

import (
	"encoding/hex"
	"strings"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/crypto"
)

// RecoverAddressFromPersonalSignature 从 personal_sign 风格签名中恢复钱包地址。
//
// 钱包登录时，前端通常使用 personal_sign 对完整 message 签名。
// Ethereum personal_sign 实际签名的不是原始 message，
// 而是加上 "\x19Ethereum Signed Message:\n" 前缀后的 hash。
// go-ethereum 的 accounts.TextHash 会生成这个 hash。
func RecoverAddressFromPersonalSignature(message string, signature string) (string, error) {
	signature = strings.TrimSpace(signature)
	signature = strings.TrimPrefix(signature, "0x")

	sig, err := hex.DecodeString(signature)
	if err != nil {
		return "", ErrInvalidSignature
	}

	if len(sig) != 65 {
		return "", ErrInvalidSignature
	}

	// MetaMask 等钱包返回的 v 通常是 27/28，
	// go-ethereum crypto.SigToPub 需要 0/1。
	if sig[64] >= 27 {
		sig[64] -= 27
	}

	if sig[64] != 0 && sig[64] != 1 {
		return "", ErrInvalidSignature
	}

	hash := accounts.TextHash([]byte(message))

	pubKey, err := crypto.SigToPub(hash, sig)
	if err != nil {
		return "", ErrInvalidSignature
	}

	return crypto.PubkeyToAddress(*pubKey).Hex(), nil
}
