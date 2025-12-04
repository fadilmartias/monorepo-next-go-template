package models

import (
	"database/sql"
	"database/sql/driver"
	"math/rand"
	"time"

	"github.com/bytedance/sonic"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID                 string         `gorm:"primaryKey;size:7" json:"id"`
	GoogleID           NullString     `gorm:"size:100"`
	DiscordID          NullString     `gorm:"size:100"`
	FacebookID         NullString     `gorm:"size:100"`
	SteamID            NullString     `gorm:"size:100"`
	AppleID            NullString     `gorm:"size:100"`
	TwitchID           NullString     `gorm:"size:100"`
	TenantID           string         `gorm:"not null;size:7;index" json:"tenant_id"`
	TotalExp           int            `gorm:"not null;default:0" json:"total_exp"`
	Exp                int            `gorm:"not null;default:0" json:"exp"`
	Level              int            `gorm:"not null;default:1" json:"level"`
	Name               string         `gorm:"not null;size:100;index" faker:"name" json:"name"`
	Email              string         `gorm:"uniqueIndex;not null;size:100" faker:"email" json:"email"`
	Phone              string         `gorm:"size:15" json:"phone"`
	Password           string         `gorm:"not null;size:100" faker:"password" json:"password"`
	Role               string         `gorm:"type:enum('superadmin','admin','user');default:'user';not null;index" json:"role"`
	RefreshToken       *string        `gorm:"type:text" json:"refresh_token"`
	TOTPSecret         *string        `gorm:"type:text" json:"totp_secret"`
	EmailVerifiedAt    *time.Time     `json:"email_verified_at"`
	TotalSpent         float64        `gorm:"not null;default:0" json:"total_spent"`
	ActiveAvatarID     NullString     `json:"active_avatar_id"`
	ActiveBackgroundID NullString     `json:"active_background_id"`
	ActiveBorderID     NullString     `json:"active_border_id"`
	ActiveTitleID      NullString     `json:"active_title_id"`
	CreatedAt          time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func GenerateID(length int) string {
	const shortIDChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = shortIDChars[seededRand.Intn(len(shortIDChars))]
	}
	return string(b)
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		u.ID = GenerateID(7)
	}

	return
}

// HashPassword mengenkripsi password sebelum disimpan
func (u *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	u.Password = string(bytes)
	return nil
}

// CheckPassword memverifikasi password
func (u *User) CheckPassword(providedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(providedPassword))
}

// NullString type untuk handle sql.NullString + JSON
type NullString struct {
	sql.NullString
}

// Implement Scanner (untuk baca dari DB)
func (ns *NullString) Scan(value any) error {
	return ns.NullString.Scan(value)
}

// Implement Valuer (untuk simpan ke DB)
func (ns NullString) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return ns.String, nil
}

// Implement MarshalJSON (untuk output ke JSON)
func (ns NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return sonic.Marshal("") // kalau NULL -> string kosong
	}
	return sonic.Marshal(ns.String)
}

func NewNullString(s string) NullString {
	return NullString{NullString: sql.NullString{String: s, Valid: s != ""}}
}

type NullInt64 struct {
	sql.NullInt64
}

func (ni *NullInt64) Scan(value any) error {
	return ni.NullInt64.Scan(value)
}

func (ni NullInt64) Value() (driver.Value, error) {
	if !ni.Valid {
		return nil, nil
	}
	return ni.Int64, nil
}

func (ni NullInt64) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return sonic.Marshal(0)
	}
	return sonic.Marshal(ni.Int64)
}

func NewNullInt64(i int64) NullInt64 {
	return NullInt64{NullInt64: sql.NullInt64{Int64: i, Valid: true}}
}

type NullFloat64 struct {
	sql.NullFloat64
}

func (nf *NullFloat64) Scan(value any) error {
	return nf.NullFloat64.Scan(value)
}

func (nf NullFloat64) Value() (driver.Value, error) {
	if !nf.Valid {
		return nil, nil
	}
	return nf.Float64, nil
}

func (nf NullFloat64) MarshalJSON() ([]byte, error) {
	if !nf.Valid {
		return sonic.Marshal(0.0)
	}
	return sonic.Marshal(nf.Float64)
}

func NewNullFloat64(f float64) NullFloat64 {
	return NullFloat64{NullFloat64: sql.NullFloat64{Float64: f, Valid: true}}
}

type NullBool struct {
	sql.NullBool
}

func (nb *NullBool) Scan(value any) error {
	return nb.NullBool.Scan(value)
}

func (nb NullBool) Value() (driver.Value, error) {
	if !nb.Valid {
		return nil, nil
	}
	return nb.Bool, nil
}

func (nb NullBool) MarshalJSON() ([]byte, error) {
	if !nb.Valid {
		return sonic.Marshal(false)
	}
	return sonic.Marshal(nb.Bool)
}

func NewNullBool(b bool) NullBool {
	return NullBool{NullBool: sql.NullBool{Bool: b, Valid: true}}
}

type NullTime struct {
	sql.NullTime
}

func (nt *NullTime) Scan(value any) error {
	return nt.NullTime.Scan(value)
}

func (nt NullTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}

func (nt NullTime) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return sonic.Marshal(nil)
	}
	return sonic.Marshal(nt.Time)
}

func NewNullTime(t time.Time) NullTime {
	return NullTime{NullTime: sql.NullTime{Time: t, Valid: true}}
}
