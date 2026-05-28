package models

type Setting struct {
	BaseModel
	Key         string     `gorm:"not null;uniqueIndex;size:255" json:"key"`
	Value       string     `gorm:"not null;type:text" json:"value"`
	Type        SettingKey `gorm:"not null;size:20;default:text;size:255" json:"type"`
	Description string     `gorm:"type:text" json:"description"`
	IsActive    bool       `gorm:"not null;default:true;index" json:"is_active"`
}
