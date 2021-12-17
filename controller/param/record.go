package param

type UpdateRecordRequest struct {
	Status int `json:"status"`
}

type RecordResponse struct {
	CardID     string `json:"card_id"`
	LastStudy  int64  `json:"last_study"`
	StudyTimes int    `json:"study_times"`
	Status     int    `json:"status"`
}

type GetRecordsRespons struct {
	Total   int              `json:"total"`
	Records []RecordResponse `json:"records"`
}
