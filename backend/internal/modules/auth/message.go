package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

const DefaultSIWEVersion = "1"

// BuildSIWEMessage 构造 SIWE 风格的钱包登录消息
//
// 这里没有让前端自己拼 message，是因为登录挑战必须由服务端生成并保存。
// Verify 阶段会根据完整 message 计算 message_hash，并从 auth_nonces 表找回挑战。
// 如果允许前端随意拼 message，就无法保证 nonce、domain、chain_id、expires_at 等字段可信。
func BuildSIWEMessage(domain, address, statement, uri string, chainID int64, nonce string, issuedAt, expiresAt time.Time) string {
	issuedAt = issuedAt.UTC()
	expiresAt = expiresAt.UTC()

	var b strings.Builder

	b.WriteString(domain)
	b.WriteString(" wants you to sign in with your Ethereum account:\n")
	b.WriteString(address)
	b.WriteString("\n\n")

	if statement != "" {
		b.WriteString(statement)
		b.WriteString("\n\n")
	}

	b.WriteString("URI: ")
	b.WriteString(uri)
	b.WriteString("\n")

	b.WriteString("Version: ")
	b.WriteString(DefaultSIWEVersion)
	b.WriteString("\n")

	b.WriteString("Chain ID: ")
	b.WriteString(fmt.Sprintf("%d", chainID))
	b.WriteString("\n")

	b.WriteString("Nonce: ")
	b.WriteString(nonce)
	b.WriteString("\n")

	b.WriteString("Issued At: ")
	b.WriteString(issuedAt.Format(time.RFC3339))
	b.WriteString("\n")

	b.WriteString("Expiration Time: ")
	b.WriteString(expiresAt.Format(time.RFC3339))

	return b.String()
}

// HashMessage 返回 message 的 sha256 hex 字符串。
//
// auth_nonces.message 是 TEXT，不适合直接做唯一索引和高效查询。
// 因此 Verify 阶段用 sha256(message) 得到固定长度的 message_hash，
// 再通过 message_hash 查找登录挑战。
func HashMessage(message string) string {
	sum := sha256.Sum256([]byte(message))
	return hex.EncodeToString(sum[:])
}
