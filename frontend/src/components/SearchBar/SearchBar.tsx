import React, { useState } from 'react';
import './SearchBar.css';

const SearchBar: React.FC = () => {
    const [searchQuery, setSearchQuery] = useState<string>('');

    const handleSearch = (e: React.ChangeEvent<HTMLInputElement>) => {
        setSearchQuery(e.target.value);
        // TODO: Implement search functionality
    };

    return (
        <div className="search-bar">
            <span className="search-icon">🔍</span>
            <input
                type="text"
                placeholder="Искать клиента, сотрудника, документ"
                className="search-input"
                value={searchQuery}
                onChange={handleSearch}
            />
        </div>
    );
};

export default SearchBar;

