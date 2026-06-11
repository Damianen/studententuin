import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import { BrowserRouter } from 'react-router';
import '@fontsource-variable/inter';
import '@fontsource-variable/fraunces';
import '@fontsource-variable/fraunces/wght-italic.css';
import '@fontsource-variable/jetbrains-mono';
import '@/styles/globals.css';
import App from '@/app';

createRoot(document.getElementById('root')!).render(
	<StrictMode>
		<BrowserRouter>
			<App />
		</BrowserRouter>
	</StrictMode>
);
