package post

type Post struct {
	UserId    int      `json:"user_id" binding:"required"`
	Total     float64  `json:"total" binding:"required"`
	Title     string   `json:"title" binding:"required"`
	Metadata  Metadata `json:"meta" binding:"required"`
	Completed bool     `json:"completed"` // bug in Gin, does not allow required here for false value!
}

type Metadata struct {
	Logins       []Login      `json:"logins" binding:"required"`
	PhoneNumbers PhoneNumbers `json:"phone_numbers" binding:"required"`
}

type Login struct {
	Time string `json:"time" binding:"required"`
	IP   string `json:"ip" binding:"required"`
}

type PhoneNumbers struct {
	Home   string `json:"home" binding:"required"`
	Mobile string `json:"mobile" binding:"required"`
}