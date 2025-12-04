import type { CreateUserInput, UpdateUserInput, User } from '../entities';

export interface IUserRepo {
	create(user: CreateUserInput): Promise<User | null>;
	update(user: UpdateUserInput): Promise<User | null>;
	delete(id: string): Promise<void>;
	findById(id: string): Promise<User | null>;
	findByEmail(email: string): Promise<User | null>;
}
