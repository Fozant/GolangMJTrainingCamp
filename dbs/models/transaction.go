package models

type Transaction struct {
	IDTransaction     uint          `gorm:"primaryKey;autoIncrement;type:bigint unsigned" json:"id"`
	MembershipID      *uint         `gorm:"type:bigint unsigned" json:"id"`
	Membership        *Membership   `gorm:"foreignKey:MembershipID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	VisitPackageID    *uint         `gorm:"column:Visit_id"`
	VisitPackage      *VisitPackage `gorm:"foreignKey:VisitPackageID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	PaymentType       string        `gorm:"column:paymentType"`
	PaymentMethod     string        `gorm:"column:paymentMethod"`
	PaymentStatus     string        `gorm:"column:paymentStatus"`
	PaymentStatusNote string        `gorm:"column:paymentStatusNote"`
	TransactionPrice  uint          `gorm:"column:transactionPrice"`
	//BuktiTransfer     []byte `gorm:"column:bukti_transfer;type:blob"`
}
