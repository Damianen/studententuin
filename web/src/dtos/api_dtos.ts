export interface ApiResponse<T> {
	code: number;
	message: string;
	// Omitted by the API when there is no payload (mutations, empty lists).
	data?: T;
}
