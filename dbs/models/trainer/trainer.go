package trainer

type Trainer struct {
	IDTrainer          uint   `gorm:"primaryKey;autoIncrement" json:"id_trainer"`
	TrainerName        string `gorm:"type:varchar(100);not null" json:"trainer_name"`
	TrainerDescription string `gorm:"type:text" json:"trainer_description"`
}
