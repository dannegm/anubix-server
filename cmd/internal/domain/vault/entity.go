package vault

type Vault struct {
	ID                string `json:"id"`
	UserID            string `json:"-"`
	Name              string `json:"name"`
	EncryptedVaultKey string `json:"encrypted_vault_key"`
	VaultKeyIV        string `json:"vault_key_iv"`
	VaultKeyAuthTag   string `json:"vault_key_auth_tag"`
}
