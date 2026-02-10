import {
	createContext,
	ReactNode,
	useContext,
	useEffect,
	useState,
} from 'react';
import type { UserDto } from '../dtos/user_dtos';
import UserController from '../controllers/user_controller';
import type { LoginUserDto } from '../dtos/auth_dtos';
import AuthController from '../controllers/auth_controller';

interface AuthContextType {
	user: UserDto | null;
	loading: boolean;
	isAuthenticated: boolean;
	login: (credentials: LoginUserDto) => Promise<void>;
	logout: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
	const [user, setUser] = useState<UserDto | null>(null);
	const [loading, setLoading] = useState(true);

	useEffect(() => {
		getUser();
	}, []);

	const getUser = async () => {
		UserController.get()
			.then(setUser)
			.catch(() => setUser(null))
			.finally(() => setLoading(false));
	};

	const login = async (credentials: LoginUserDto) => {
		await AuthController.login(credentials);
		getUser();
	};

	const logout = async () => {
		await AuthController.logout();
		setUser(null);
		setLoading(false);
	};

	return (
		<AuthContext.Provider
			value={{ user, loading, isAuthenticated: !!user, login, logout }}
		>
			{children}
		</AuthContext.Provider>
	);
}

export function useAuth() {
	const context = useContext(AuthContext);
	if (context === undefined) {
		throw new Error('useAuth can only be used within the authprovider');
	}
	return context;
}
