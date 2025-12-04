export interface IEncryptionService {
	encrypt(string: string, key: string): Promise<string>;
	decrypt(string: string, key: string): Promise<string>;
}
