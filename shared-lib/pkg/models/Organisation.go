
package models

type Organisation struct {

	ID int `json:"id"`
	Name string `json:"name"`
	Url string `json:"url"`
	ContactName string `json:"contact_name"`
	ContactEmail string `json:"contact_email"`
	ContactPhone string `json:"contact_phone"`
	LogoAttachmentId int `json:"logo_attachment_id"`
	SmallLogoAttachmentId int `json:"small_logo_attachment_id"`
}
