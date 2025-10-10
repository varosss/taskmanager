package item

type Task struct {
	Id       int    `json:"id"`
	Title    string `json:"title"`
	Status   Status `json:"status"`
	Category string `json:"category"`
}
