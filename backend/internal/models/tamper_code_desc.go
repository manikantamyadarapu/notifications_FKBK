package models

type TamperCodeDesc struct {
	TamperCode int    `json:"tamper_code" gorm:"column:tamper_code;primaryKey"`
	TamperDesc string `json:"tamper_desc" gorm:"column:tamper_desc"`
}

func (TamperCodeDesc) TableName() string { return "tamper_code_desc" }

