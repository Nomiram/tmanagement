package headers

var CONNSTR = "user=postgres password=qwerty dbname=VS sslmode=disable"

// var CONNSTR = "host=db port=5432 user=postgres password=postgres sslmode=disable"

var CONNSTRWDB = CONNSTR + " dbname=vs"

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

func AddDBinCONNSTR(db string) {
	CONNSTR = "host=" + db + " port=5432 user=postgres password=postgres sslmode=disable"
	CONNSTRWDB = CONNSTR + " dbname=VS"
}

// type redisapi
