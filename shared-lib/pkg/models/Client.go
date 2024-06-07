
package models

type Client struct {

	ID int `json:"id"`
	CreatorUid int `json:"creator_uid"`
	Name string `json:"name"`
	Address string `json:"address"`
	Url string `json:"url"`
	LogoAttachmentId int `json:"logo_attachment_id"`
	SmallLogoAttachmentId int `json:"small_logo_attachment_id"`
}
