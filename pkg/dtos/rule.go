// Copyright (C) 2023-2024 IOTech Ltd

package dtos

type Rule struct {
	Name string `json:"name" validate:"required"`
	Rule []byte `json:"rule" validate:"required"`
}

func NewRule(name string, rule []byte) Rule {
	return Rule{
		Name: name,
		Rule: rule,
	}
}
