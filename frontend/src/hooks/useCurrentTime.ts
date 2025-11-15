import { useState, useEffect, useRef } from 'react';

export const useCurrentTime = (): string => {
    const [currentTime, setCurrentTime] = useState<string>(() => {
        const now = new Date();
        const hours = now.getHours().toString().padStart(2, '0');
        const minutes = now.getMinutes().toString().padStart(2, '0');
        return `${hours}:${minutes}`;
    });

    const intervalRef = useRef<NodeJS.Timeout | null>(null);

    useEffect(() => {
        const updateTime = () => {
            const now = new Date();
            const hours = now.getHours().toString().padStart(2, '0');
            const minutes = now.getMinutes().toString().padStart(2, '0');
            setCurrentTime(`${hours}:${minutes}`);
        };

        // Обновляем сразу при монтировании
        updateTime();
        
        // Устанавливаем интервал
        intervalRef.current = setInterval(updateTime, 1000);

        return () => {
            if (intervalRef.current) {
                clearInterval(intervalRef.current);
            }
        };
    }, []);

    return currentTime;
};
