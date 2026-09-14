export const API_BASE_URL =
    import.meta.env.VITE_API_BASE_URL ?? 'https://urlvio.onrender.com/api/v1'

export class ApiError extends Error {
    readonly status: number

    constructor(message: string, status: number) {
        super(message)
        this.name = 'ApiError'
        this.status = status
    }
}

interface ApiFetchOptions {
    method?: 'GET' | 'POST' | 'PUT' | 'DELETE'
    body?: unknown
    signal?: AbortSignal
}

async function extractErrorMessage(response: Response): Promise<string> {
    try {
        const data: unknown = await response.json()
        if (
        data &&
        typeof data === 'object' &&
        'error' in data &&
        typeof (data as { error: unknown }).error === 'string'
        ) {
        return (data as { error: string }).error
        }
    } catch {
        // Body wasn't JSON - fall through to the status text.
    }

    return response.statusText || `Request failed with status ${response.status}`
}

export async function apiFetch<T>(
    path: string,
    token: string | null,
    options: ApiFetchOptions = {}
    ): Promise<T> {
    const headers = new Headers({ Accept: 'application/json' })

    if (options.body !== undefined) {
        headers.set('Content-Type', 'application/json')
    }

    if (token) {
        headers.set('Authorization', `Bearer ${token}`)
    }

    let response: Response

    try {
        response = await fetch(`${API_BASE_URL}${path}`, {
        method: options.method ?? 'GET',
        headers,
        body: options.body !== undefined ? JSON.stringify(options.body) : undefined,
        signal: options.signal,
        })
    } catch {
        throw new ApiError('Unable to reach the server. Check your connection and try again.', 0)
    }

    if (!response.ok) {
        throw new ApiError(await extractErrorMessage(response), response.status)
    }

    if (response.status === 204) {
        return undefined as T
    }

    return (await response.json()) as T
}
