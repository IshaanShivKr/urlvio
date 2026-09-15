export function validateUrl(value: string): string | null {
    const trimmed = value.trim()

    if (!trimmed) {
        return 'Enter a destination URL.'
    }

    let parsed: URL
    try {
        parsed = new URL(trimmed)
    } catch {
        return 'Enter a valid URL.'
    }

    if (parsed.protocol !== 'http:' && parsed.protocol !== 'https:') {
        return 'Enter a valid URL.'
    }

    return null
}
