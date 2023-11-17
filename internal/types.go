package internal

type Language struct {
	Name       string
	Sphere     string `json:"type"`
	Extensions []string
}

type Statistic struct {
	Commits int
	Lines   int
	Files   int
}
