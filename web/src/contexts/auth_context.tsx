import {
	createContext,
	ReactNode,
	useContext,
	useEffect,
	useState,
} from 'react';

interface AuthContextType {
	user: any | null;
	loading: boolean;
	isAuthenticated: boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
	const [user, setUser] = useState<any | null>(null);
	const [loading, setLoading] = useState(true);

	useEffect(() => {
		//TODO get user and set it to user
	}, []);

	return (
		<AuthContext.Provider
			value={{ user, loading, isAuthenticated: !!user }}
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
