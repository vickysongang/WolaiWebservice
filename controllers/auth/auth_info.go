package auth

type authInfo struct {
	Id          int64  `json:"id"`
	Nickname    string `json:"nickname"`
	Avatar      string `json:"avatar"`
	Gender      int64  `json:"gender"`
	AccessRight int64  `json:"accessRight"`
	Token       string `json:"token"`
}
