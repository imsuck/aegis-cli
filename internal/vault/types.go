package vault

// Vault represents the top-level Aegis vault structure
type Vault struct {
	Version int    `json:"version"`
	Header  Header `json:"header"`
	DB      string `json:"db"` // Base64 encoded when encrypted
}

// Header contains encryption parameters and key slots
type Header struct {
	Slots  []Slot `json:"slots"`
	Params Params `json:"params"`
}

// Slot represents a credential slot for vault decryption
type Slot struct {
	Type      int       `json:"type"`
	UUID      string    `json:"uuid"`
	Key       string    `json:"key"` // Hex encoded
	KeyParams KeyParams `json:"key_params"`
	N         int       `json:"n"` // scrypt N parameter
	R         int       `json:"r"` // scrypt r parameter
	P         int       `json:"p"` // scrypt p parameter
	Salt      string    `json:"salt"` // Hex encoded
}

// KeyParams holds AES-GCM nonce and tag
type KeyParams struct {
	Nonce string `json:"nonce"` // Hex encoded
	Tag   string `json:"tag"`   // Hex encoded
}

// Params holds vault-level encryption parameters
type Params struct {
	Nonce string `json:"nonce"` // Hex encoded
	Tag   string `json:"tag"`   // Hex encoded
}

// Content represents the decrypted vault contents
type Content struct {
	Version int     `json:"version"`
	Entries []Entry `json:"entries"`
	Groups  []Group `json:"groups"`
}

// Entry represents a single OTP entry
type Entry struct {
	Type     string   `json:"type"`
	UUID     string   `json:"uuid"`
	Name     string   `json:"name"`
	Issuer   string   `json:"issuer"`
	Note     string   `json:"note"`
	Icon     *string  `json:"icon"`
	IconMime *string  `json:"icon_mime"`
	IconHash *string  `json:"icon_hash"`
	Favorite bool     `json:"favorite"`
	Info     Info     `json:"info"`
	Groups   []string `json:"groups"`
}

// Info holds OTP-specific information
type Info struct {
	Secret  string `json:"secret"`
	Algo    string `json:"algo"`
	Digits  int    `json:"digits"`
	Period  int    `json:"period,omitempty"`
	Counter uint64 `json:"counter,omitempty"`
}

// Group represents an entry group
type Group struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}
