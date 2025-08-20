package hhclient

type HHClient interface {
	GetToken() (string, error)
	SearchVacancies(token string, keywords []string) ([]Vacancy, error)
	ApplyVacancy(token, vacancyID, resumeID string) error
}

type Repo interface {
	AlreadyResponded(vacancyID, resumeID string) bool
	SaveResponse(vacancyID, resumeID string) error
}
