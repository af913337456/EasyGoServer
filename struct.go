package main

// nullTag: 0 代表可以在读取post过来的参数时候为 null. 1 不能为 null 
type LUser struct { 
	Id			int	`json:"id" nullTag:"1"`
	U_user_id			*string	`json:"u_user_id" nullTag:"1"`
	U_time			*string	`json:"u_time" nullTag:"1"`
	U_open_time			int	`json:"u_open_time" nullTag:"1"`
}
type LComment struct { 
	Id			int	`json:"id" nullTag:"1"`
	U_user_id			*string	`json:"u_user_id" nullTag:"1"`
	U_time			*string	`json:"u_time" nullTag:"1"`
	U_content			*string	`json:"u_content" nullTag:"1"`
}
