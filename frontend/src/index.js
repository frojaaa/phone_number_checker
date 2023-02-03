import React, {createContext} from 'react';
import ReactDOM from 'react-dom/client';
import App from './App';
import './index.css';

export const Context = createContext(null);

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
    <React.StrictMode>
        <App/>
    </React.StrictMode>
);