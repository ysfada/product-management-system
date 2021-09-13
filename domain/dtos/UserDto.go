package dtos

type UserDto struct {
	ID          int    `json:"id"`
	Username    string `json:"username"`
	IsActive    bool   `json:"is_active"`
	IsStaff     bool   `json:"is_staff"`
	IsSuperuser bool   `json:"is_superuser"`
}
