import { BrowserRouter, Navigate, Route, Routes } from 'react-router-dom';
import { AuthProvider } from '@/contexts/auth_context';
import { ThemeProvider } from '@/components/theme-provider';
import { useAuth } from '@/contexts/auth_context';

function App() {
	return (
		<BrowserRouter>
				<AuthProvider>
					<Routes>
						<Route path="/" element={<div></div>} />
					</Routes>
				</AuthProvider>
		</BrowserRouter>
	);
}

export default App;
