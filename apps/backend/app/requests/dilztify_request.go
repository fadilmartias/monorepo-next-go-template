package requests

type DilztifySendMessageRequest struct {
	Phone   string `gorm:"type:text;not null" json:"phone"` // required
	Message string `gorm:"type:text" json:"message"`        // required
}
