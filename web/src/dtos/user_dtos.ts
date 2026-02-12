export interface UserDto {
	email: string;
	name: string;
	status: string;
}

export interface RegisterUserDto {
	email: string;
	password: string;
	name: string;
}

export interface UpdateUserDto {
	email?: string;
	name?: string;
	status?: string;
}
