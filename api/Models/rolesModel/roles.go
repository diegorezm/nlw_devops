package rolesmodel

type Roles struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type UserRoles struct {
	RoleId int `json:"role_id"`
	UserId int `json:"user_id"`
}

func NewUserRoles(roleId, userId int) *UserRoles {
	return &UserRoles{
		RoleId: roleId,
		UserId: userId,
	}
}

func NewRoles(id int, name string) *Roles {
	return &Roles{
		Id:   id,
		Name: name,
	}
}
