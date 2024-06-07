
package models

type User struct {

	ID int `json:"id"`
	SubjectId string `json:"subject_id"`
	Active bool `json:"active"`
	FullName string `json:"full_name"`
	ShortBio string `json:"short_bio"`
	Username string `json:"username"`
	Email string `json:"email"`
	Role string `json:"role"`
}
