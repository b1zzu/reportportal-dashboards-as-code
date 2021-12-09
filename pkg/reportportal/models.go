package reportportal

type EntryCreated struct {
	ID int `json:"id"`
}

type OperationCompletion struct {
	Message string `json:"message"`
}
