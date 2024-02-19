package response

import "GinLearning/gin_application/model"

type UserDto struct {
	Name      string `json:"name"`
	Telephone string `json:"telephone"`
}

// ToUserDto 将User对象转换为UserDto对象 DTO: Data Transfer Object 数据传输对象，用于展示层与服务层之间的数据传输对象
func ToUserDto(user model.User) UserDto {
	return UserDto{
		Name:      user.Name,
		Telephone: user.Telephone,
	}
}
