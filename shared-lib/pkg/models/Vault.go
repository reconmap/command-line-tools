
package models

type Vault struct {

	ID int `json:"id"`
	InsertTs string `json:"insert_ts"`
	UpdateTs string `json:"update_ts"`
	Name string `json:"name"`
	Value string `json:"value"`
	Reportable bool `json:"reportable"`
	Note string `json:"note"`
	Type string `json:"type"`
	ProjectId int `json:"project_id"`
}
