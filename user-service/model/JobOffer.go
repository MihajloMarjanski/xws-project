package model

type JobOffer struct {
	ID             uint   `json:"id"`
	JobPosition    string `json:"jobPosition"`
	CompanyName    string `json:"companyName"`
	JobInfo        string `json:"jobInfo"`
	Qualifications string `json:"qualifications"`
	ApiKey         string `json:"apiKey"`
}
