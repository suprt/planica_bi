import React from 'react';
import './PageTemplate.css';

interface PageTemplateProps {
    title: string;
    children?: React.ReactNode;
}

const PageTemplate: React.FC<PageTemplateProps> = ({ title, children }) => {
    return (
        <div className="page-template">
            <h1 className="page-title">{title}</h1>
            {children || (
                <div className="page-content">
                    <p>Раздел находится в разработке</p>
                </div>
            )}
        </div>
    );
};

export default PageTemplate;

