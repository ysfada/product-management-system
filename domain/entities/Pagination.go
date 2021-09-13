package entities

type Pagination struct {
	TotalPage    int `json:"total_page"`
	CurrentPage  int `json:"current_page"`
	NextPage     int `json:"next_page"`
	PreviousPage int `json:"previous_page"`
	Count        int `json:"count"`
	Size         int `json:"size"`
}
