package response

type LoginRequestJWT struct {
	Username string `json:"username" example:"edoardo"`
	Password string `json:"password" example:"Nethesis,1234"`
}

type LoginResponseJWT struct {
	Code   int    `json:"code" example:"200"`
	Token  string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY3Rpb25zIjpbXSwiZXhwIjoxNjE5NTM0OTQ4LCJpZCI6ImVkb2FyZG8iLCJvcmlnX2lhdCI6MTYxODkzMDE0OCwicm9sZSI6IiJ9.bNRFa7MCQK-rTczOjLveXEWBqhjK-FWhnUPD3_ixcCI"`
	Expire string `json:"expire" example:"2021-04-27T16:49:08+02:00"`
}

type StatusOK struct {
	Code    int         `json:"code" example:"200"`
	Message string      `json:"code" example:"Success"`
	Data    interface{} `json:"code"`
}

type StatusCreated struct {
	Code    int         `json:"code" example:"201"`
	Message string      `json:"message" example:"Created"`
	Data    interface{} `json:"data"`
}

type StatusBadRequest struct {
	Code    int         `json:"code" example:"400"`
	Message string      `json:"message" example:"Bad request"`
	Data    interface{} `json:"data"`
}

type StatusUnauthorized struct {
	Code    int         `json:"code" example:"401"`
	Message string      `json:"message" example:"Unauthorized"`
	Data    interface{} `json:"data"`
}

type StatusForbidden struct {
	Code    int         `json:"code" example:"403"`
	Message string      `json:"message" example:"Forbidden"`
	Data    interface{} `json:"data"`
}

type StatusNotFound struct {
	Code    int         `json:"code" example:"404"`
	Message string      `json:"message" example:"Not found"`
	Data    interface{} `json:"data"`
}

type StatusInternalServerError struct {
	Code    int         `json:"code" example:"500"`
	Message string      `json:"message" example:"Internal server error"`
	Data    interface{} `json:"data"`
}
