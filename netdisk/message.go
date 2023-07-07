package netdisk

type Resp struct {
	Errno     int    `json:"errno"`
	Path      string `json:"path"`
	RequestId uint64 `json:"request_id"`
}
