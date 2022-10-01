package handlers

type Order struct {
	Order_name string `json:"order_name"`
	Start_date string `json:"start_date"`
}

type Delorder struct {
	Order_name string `json:"order_name"`
}

type Task struct {
	Task       string `json:"task"`
	Order_name string `json:"order_name"`
	Duration   int    `json:"duration"`
	Resource   int    `json:"resource"`
	Pred       string `json:"pred"`
}

type TaskEn struct {
	Task       string   `json:"task"`
	Order_name string   `json:"order_name"`
	Duration   int      `json:"duration"`
	Resource   int      `json:"resource"`
	Pred       []string `json:"pred"`
}

type TaskDel struct {
	Task       string `json:"task"`
	Order_name string `json:"order_name"`
}
