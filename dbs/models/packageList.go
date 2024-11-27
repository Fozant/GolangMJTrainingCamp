package models

type PackageList struct {
	IDPackage   uint   `gorm:"primaryKey;autoIncrement;column:id_package"`
	PackageName string `gorm:"column:package_name;size:255;not null"`
	Price       uint   `gorm:"column:price;not null"`
	Duration    *uint  `gorm:"column:duration""`
	Status      string `gorm:"column:status;size:50;not null"`
	Type        string `gorm:"column:type;size:50;not null"`
	VisitNumber *uint  `gorm:"column:visit_number""`
}
