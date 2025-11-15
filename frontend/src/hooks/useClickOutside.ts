import { useEffect, RefObject, useRef } from 'react';

export const useClickOutside = (
    ref: RefObject<HTMLElement | null>,
    handler: () => void
) => {
    const handlerRef = useRef(handler);

    // Обновляем ref при каждом изменении handler
    useEffect(() => {
        handlerRef.current = handler;
    }, [handler]);

    useEffect(() => {
        const handleClickOutside = (event: MouseEvent) => {
            if (ref.current && !ref.current.contains(event.target as Node)) {
                handlerRef.current();
            }
        };

        // Используем capture phase для лучшей производительности
        document.addEventListener('mousedown', handleClickOutside, true);
        return () => {
            document.removeEventListener('mousedown', handleClickOutside, true);
        };
    }, [ref]);
};
