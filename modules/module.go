package modules

type Module struct {
	Name      string `json:"name"`
	Form      string `json:"form"`
	URL       string `json:"url"`
	Rev       string `json:"rev"`
	Dynamic   bool   `json:"dynamic"`
	Shprov    string `json:"shprov"`
	ShprovDir string `json:"shprovdir"`
}
