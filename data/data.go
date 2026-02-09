package data

import (
	_ "embed"
	"encoding/json"
	"html/template"
)

//go:embed site.json
var siteJSON []byte

type Site struct {
	Name     string `json:"name"`
	Role     string `json:"role"`
	Tagline  string `json:"tagline"`
	Email    string `json:"email"`
	Github   string `json:"github"`
	Telegram string `json:"telegram"`
	Linkedin string `json:"linkedin"`

	About     About    `json:"about"`
	Skills    []Skill  `json:"skills"`
	Social    []Link   `json:"social"`
	Interests []string `json:"interests"`
}

type About struct {
	Intro   template.HTML `json:"intro"`
	Mission template.HTML `json:"mission"`
}

type Skill struct {
	Name  string `json:"name"`
	Level string `json:"level,omitempty"`
	Icon  string `json:"icon,omitempty"`
}

type Link struct {
	Name string `json:"name"`
	Url  string `json:"url"`
	Icon string `json:"icon,omitempty"`
}

var Data Site

func init() {
	err := json.Unmarshal(siteJSON, &Data)
	if err != nil {
		panic("cannot parse site.json: " + err.Error())
	}
}
