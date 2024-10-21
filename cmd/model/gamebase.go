package model

type GameBase struct {
	ID        int    `json:"id"`
	RoleCount int    `json:"role_count"`
	RolesList []Role `json:"roles_list"`
}
type Role struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Type int    `json:"type"`
}
