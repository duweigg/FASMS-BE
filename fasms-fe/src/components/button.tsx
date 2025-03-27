'use client';

import { useRouter } from 'next/navigation';

interface NavigationButtonProps {
    text: string;
    page: string;
}

export const NavigationButton = ({ text, page }: NavigationButtonProps) => {
    const router = useRouter();
    const handleClick = () => {
        router.push(page)
    };
    return (
        <button
            onClick={handleClick}
            className={"border-2 px-5 py-2 rounded-2xl w-full"}>
            {text}
        </button>
    );
};

interface BackButtonProps {
    text: string;
}

export const BackButton = ({ text }: BackButtonProps) => {
    const router = useRouter();
    const handleClick = () => {
        router.back()
    };
    return (
        <button
            onClick={handleClick}
            className={"px-10 py-3 border-0"}>
            {text}
        </button>
    );
};