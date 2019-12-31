package resp

// UserInfoResp 用户信息解析
type UserInfoResp struct {
	Id        string  `json:"id"`         // id
	Nickname  string  `json:"nickname"`   // nickname
	Email     string  `json:"email"`      // email
	Rcode     string  `json:"rcode"`      // rcode
	FatherId  string  `json:"father_id"`  // father_id
	HeadImage string  `json:"head_image"` // head_image
	Fans      float64 `json:"fans"`       // fans
}
