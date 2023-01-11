package db

type User struct {
	Id          string
	Name        string
	StudentName string
	Password    string
	Email       string
	Salt        string
	RoleIds     string
	Enabled 	bool
}

func (User) TableName() string {
	return "t_user"
}