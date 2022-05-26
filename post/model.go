package post

type Post struct {
	UserId    int      `json:"user_id"`
	Total     float64  `json:"total"`
	Title     string   `json:"title"`
	Metadata  Metadata `json:"meta"`
	Completed bool     `json:"completed"`
}

type Metadata struct {
	Logins       []Login      `json:"logins"`
	PhoneNumbers PhoneNumbers `json:"phone_numbers"`
}

type Login struct {
	Time string `json:"time"`
	IP   string `json:"ip"`
}

type PhoneNumbers struct {
	Home   string `json:"home"`
	Mobile string `json:"mobile"`
}