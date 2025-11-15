import React from 'react';
import './Home.css';

const Home: React.FC = () => {
    return (
        <div className="home-page">
            <div className="content-placeholder">
                <h2>Добро пожаловать в Planica</h2>
                <p>Выберите раздел в меню для начала работы</p>
            </div>
        </div>
    );
};

export default Home;

