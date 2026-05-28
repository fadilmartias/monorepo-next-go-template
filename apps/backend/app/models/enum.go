package models

type SessionStatus string

const (
	SessionBotHandling   SessionStatus = "BOT_HANDLING"
	SessionWaitingAgent  SessionStatus = "WAITING_FOR_AGENT"
	SessionAgentHandling SessionStatus = "AGENT_HANDLING"
	SessionResolved      SessionStatus = "RESOLVED"
)

type SenderRole string

const (
	SenderUser  SenderRole = "USER"
	SenderBot   SenderRole = "BOT"
	SenderAgent SenderRole = "AGENT"
)

type UserRole string

const (
	UserRoleUser       UserRole = "user"
	UserRoleSuperAdmin UserRole = "superadmin"
	UserRoleAdmin      UserRole = "admin"
)

type ChangeType string

const (
	ChangeTypeUpgrade   ChangeType = "upgrade"
	ChangeTypeDowngrade ChangeType = "downgrade"
)

type SettingKey string

const (
	SettingKeyText    SettingKey = "text"
	SettingKeyNumber  SettingKey = "number"
	SettingKeyBoolean SettingKey = "boolean"
)

type QueuePriority string

const (
	QueuePriorityLow    QueuePriority = "low"
	QueuePriorityMedium QueuePriority = "medium"
	QueuePriorityHigh   QueuePriority = "high"
)

type CategoryType string

const (
	CategoryTypeProduct            CategoryType = "product"
	CategoryTypeTransaction        CategoryType = "transaction"
	CategoryTypePaymentMethod      CategoryType = "payment-method"
	CategoryTypeArticle            CategoryType = "article"
	CategoryTypeProductSubcategory CategoryType = "product-subcategory"
	CategoryTypeProductParent      CategoryType = "product-parent"
)

type TopupStatus string

const (
	TopupStatusPending TopupStatus = "pending"
	TopupStatusSuccess TopupStatus = "success"
	TopupStatusFailed  TopupStatus = "failed"
	TopupStatusExpired TopupStatus = "expired"
)
