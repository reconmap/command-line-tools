
package models

type CommandUsage struct {

	ID int `json:"id"`
	CommandId int `json:"command_id"`
	CreatorUid int `json:"creator_uid"`
	Name string `json:"name"`
	Description string `json:"description"`
	ExecutablePath string `json:"executable_path"`
	Arguments string `json:"arguments"`
	OutputFilename string `json:"output_filename"`
	OutputParser string `json:"output_parser"`
}
