package model

type GameBase struct {
	ID        int      `json:"id"`
	RoleCount int      `json:"role_count"`
	RolesList []string `json:"roles_list"`
}
