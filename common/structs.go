package common

type Route struct {
	Uri        string
	Controller string
	Action     []string
	Params     *Params
	Methods    []string
}

type Params struct {
	IntOnly    []string
	StringOnly []string
	Enum       map[string]string
	LimitQuery string
}
