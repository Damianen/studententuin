import React from 'react';
import ReactDOM from 'react-dom/client';

import '../public/style.css';

function App() {
    const [count, setCount] = React.useState(0);

    const increment = () => {
        setCount(count + 1);
    };

    const decrement = () => {
        setCount(count - 1);
    };

    return (
        <div>
            <h1 className='text-4xl text-blue'>React is working!!!!</h1>
            <p>Count: {count}</p>
            <button className='w-60 bg-black text-white shadow rounded' onClick={increment}>Increment</button>
            <button className='w-60 bg-black text-white shadow rounded' onClick={decrement}>Decrement</button>
        </div>
    );
};

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
    <App />,
);