package device

type Device struct {
	ID          string `json:"id"`
	UserID      string `json:"-"`
	Name        string `json:"name"`
	Fingerprint string `json:"fingerprint"`
	DeviceType  string `json:"device_type"`
}
