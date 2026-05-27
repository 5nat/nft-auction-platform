CREATE TABLE users (
   id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
   status VARCHAR(32) NOT NULL DEFAULT 'active',
   created_at DATETIME NOT NULL,
   updated_at DATETIME NOT NULL
);

CREATE TABLE wallets (
     id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
     user_id BIGINT UNSIGNED NOT NULL,
     address VARCHAR(64) NOT NULL,
     chain_id BIGINT UNSIGNED NOT NULL,
     is_primary BOOLEAN NOT NULL DEFAULT TRUE,
     created_at DATETIME NOT NULL,
     updated_at DATETIME NOT NULL,

     UNIQUE KEY uk_wallet_address_chain (address, chain_id),
     KEY idx_wallet_user_id (user_id)
);

CREATE TABLE auth_nonces (
     id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
     address VARCHAR(64) NOT NULL,
     chain_id BIGINT UNSIGNED NOT NULL,
     nonce VARCHAR(128) NOT NULL,
     message_hash CHAR(64) NOT NULL,
     message TEXT NOT NULL,
     domain VARCHAR(255) NOT NULL,
     uri VARCHAR(512) NOT NULL,
     issued_at DATETIME NOT NULL,
     expires_at DATETIME NOT NULL,
     used_at DATETIME NULL,
     created_at DATETIME NOT NULL,

     UNIQUE KEY uk_auth_nonce_nonce (nonce),
     UNIQUE KEY uk_auth_nonce_message_hash (message_hash),
     KEY idx_auth_nonce_address_chain (address, chain_id),
     KEY idx_auth_nonce_expires_at (expires_at)
);