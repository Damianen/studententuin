export type ApplicationStatus = 'deploying' | 'running' | 'stopped' | 'failed' | 'building';

export interface Application {
	id: string;
	name: string;
	userId: string;
	subdomainId: string;
	repositoryUrl?: string;
	branch?: string;
	framework?: string;
	buildCommand?: string;
	startCommand?: string;
	environmentVariables?: Record<string, string>;
	dockerImage?: string;
	dockerContainerId?: string;
	dockerContainerName?: string;
	ports?: number[];
	volumes?: string[];
	memoryLimit?: string;
	cpuLimit?: string;
	status: ApplicationStatus;
	lastDeployedAt?: Date;
	createdAt: Date;
	updatedAt: Date;
}

export type CreateApplicationInput = Omit<
	Application,
	'id' | 'createdAt' | 'updatedAt' | 'lastDeployedAt' | 'status'
>;

export type UpdateApplicationInput = Partial<
	Omit<Application, 'id' | 'userId' | 'createdAt' | 'updatedAt'>
>;
