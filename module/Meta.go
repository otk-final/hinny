package module


type ApiTag struct {
	Name        string `json:"serviceName"`
	Description string `json:"description"`
	PathCount   int    `json:"pathCount"`
}

type ApiPath struct {
	PrimaryId   string         `json:"primary_id"`
	Tag         *ApiTag        `json:"service"`
	TagName     string         `json:"tag_name"`
	Path        string         `json:"path"`
	Description string         `json:"description"`
	Method      string         `json:"method"`
	Parameters  []interface{}  `json:"parameters"`
	Definition  *ApiDefinition `json:"responses"`
	Deprecated  bool           `json:"deprecated"`
}

type ApiDefinition struct {
	Properties  map[string]interface{} `json:"properties"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
}

