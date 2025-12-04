export interface IContainerService {
	createContainer(
		image: string,
		name: string,
		memoryLimit: string,
		cpuLimit: string,
		port: number[]
	): Promise<void>;
	stopContainer(name: string): Promise<void>;
	startContainer(name: string): Promise<void>;
	deleteContainer(name: string): Promise<void>;
}
