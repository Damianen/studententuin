import {
	createContext,
	useCallback,
	useContext,
	useEffect,
	useState,
} from 'react';
import type { ReactNode } from 'react';
import { toast } from 'sonner';
import type { RegisterUserDto, UpdateUserDto, UserDto } from '../dtos/user_dtos';
import UserController from '../controllers/user_controller';
import type { LoginUserDto } from '../dtos/auth_dtos';
import AuthController from '../controllers/auth_controller';
import { setOnUnauthorized } from '../services/api_service';

interface AuthContextType {
	user: UserDto | null;
	loading: boolean;
	isAuthenticated: boolean;
	login: (credentials: LoginUserDto) => Promise<void>;
	logout: () => Promise<void>;
	register: (values: RegisterUserDto) => Promise<void>;
	updateUser: (values: UpdateUserDto) => Promise<void>;
	deleteUser: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
	const [user, setUser] = useState<UserDto | null>(null);
	const [loading, setLoading] = useState(true);

	const getUser = useCallback(async () => {
		await UserController.get()
			.then(setUser)
			.catch(() => setUser(null))
			.finally(() => setLoading(false));
	}, []);

	useEffect(() => {
		getUser();
	}, [getUser]);

	useEffect(() => {
		setOnUnauthorized(() => {
			setUser((previous) => {
				if (previous) {
					toast.error('Your session expired. Please sign in again.');
				}
				return null;
			});
		});
		return () => setOnUnauthorized(null);
	}, []);

	const login = async (credentials: LoginUserDto) => {
		await AuthController.login(credentials);
		await getUser();
	};

	const logout = async () => {
		await AuthController.logout();
		setUser(null);
		setLoading(false);
	};

	const register = async (values: RegisterUserDto) => {
		await UserController.register(values);
		await login({ email: values.email, password: values.password });
	};

	const updateUser = async (values: UpdateUserDto) => {
		await UserController.patch(values);
		await getUser();
	};

	const deleteUser = async () => {
		await UserController.delete();
		setUser(null);
	};

	return (
		<AuthContext.Provider
			value={{
				user,
				loading,
				isAuthenticated: !!user,
				login,
				logout,
				register,
				updateUser,
				deleteUser,
			}}
		>
			{children}
		</AuthContext.Provider>
	);
}

export function useAuth() {
	const context = useContext(AuthContext);
	if (context === undefined) {
		throw new Error('useAuth can only be used within AuthProvider');
	}
	return context;
}
