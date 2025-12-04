export interface User {
	id: string;
	email: string;
	name: string;
	avatarUrl?: string;
	githubId: string;
	githubUsername: string;
	githubAccessToken?: string;
	createdAt: Date;
	updatedAt: Date;
}

export type CreateUserInput = Omit<User, 'id' | 'createdAt' | 'updatedAt'>;

export type UpdateUserInput = Partial<Omit<User, 'id' | 'createdAt' | 'updatedAt'>>;
