import type { ApplicationDto } from './application_dtos';
import type { DatabaseDto } from './database_dtos';

export interface SubdomainDto {
	id: string;
	name: string;
	fullDomain: string;
	isActive: boolean;
}

export interface SubdomainListItemDto {
	id: string;
	name: string;
	fullDomain: string;
	isActive: boolean;
	database?: DatabaseDto;
	application?: ApplicationDto;
}

export interface CreateSubdomainDto {
	name: string;
	fullDomain: string;
}

export interface UpdateSubdomainDto {
	name?: string;
	fullDomain?: string;
}
