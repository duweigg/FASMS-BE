export const ApiCallFunction = async <T, TPayload = unknown>(
    method: string,
    url: string,
    payload?: TPayload,
): Promise<T | undefined> => {
    try {
        const response = await fetch(url, {
            method,
            headers: {
                'Content-Type': 'application/json',
            },
            body: payload ? JSON.stringify(payload) : undefined,
        });

        if (!response.ok) {
            if (response.status === 404) {
                console.warn(`API 404 at ${url}`);
                return undefined;
            }
            throw new Error('Network response was not ok');
        }

        const data: T = await response.json();
        return data;
    } catch (error) {
        console.error('There was a problem with the fetch operation:', error);
    }
};
