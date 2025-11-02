package models

type IPInfo struct {
	IP      string `json:"ip"`
	Country string `json:"country"`
}

func (u *IPInfo) TableName() string {
	return "ip_infos"
}
