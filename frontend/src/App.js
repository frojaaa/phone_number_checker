import React from "react";
import {BrowserRouter as Router} from "react-router-dom";
import AppRouter from "./components/AppRouter"
import 'bootstrap/dist/css/bootstrap.min.css';
import './App.css';
import Header from "./components/Header";

// import "slick-carousel/slick/slick.css";
// import "slick-carousel/slick/slick-theme.css";

function App() {
    return (
        <Router>
            <main>
                <AppRouter/>
            </main>
        </Router>
    );
}

export default App;
