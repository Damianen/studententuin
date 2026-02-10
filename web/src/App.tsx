import { BrowserRouter, Route, Routes } from 'react-router';
import { AuthProvider } from './contexts/auth_context';
import { MainLayout } from './components/layout/main_layout';
import { LoginPage } from './pages/login';
import { HomePage } from './pages/home';

function App() {
	return (
		<BrowserRouter>
			<AuthProvider>
				<MainLayout>
					<Routes>
						<Route path="/" element={<HomePage />} />
						<Route path="/login" element={<LoginPage />} />
					</Routes>
				</MainLayout>
			</AuthProvider>
		</BrowserRouter>
	);
}

export default App;
