import { useAuth } from '@clerk/clerk-react'
import { useMemo } from 'react'

import { apiFetch } from './api'
import type { CreateLinkInput, Link, LinkListItem, LinkStats, UpdateLinkInput } from '../types/link'

export function useLinksApi() {
    const { getToken } = useAuth()

    return useMemo(
        () => ({
        list: async (): Promise<LinkListItem[]> => {
            const token = await getToken()
            return apiFetch<LinkListItem[]>('/urls', token)
        },

        create: async (input: CreateLinkInput): Promise<Link> => {
            const token = await getToken()
            return apiFetch<Link>('/urls', token, { method: 'POST', body: input })
        },

        get: async (code: string): Promise<Link> => {
            const token = await getToken()
            return apiFetch<Link>(`/urls/${encodeURIComponent(code)}`, token)
        },

        update: async (code: string, input: UpdateLinkInput): Promise<Link> => {
            const token = await getToken()
            return apiFetch<Link>(`/urls/${encodeURIComponent(code)}`, token, {
            method: 'PUT',
            body: input,
            })
        },

        remove: async (code: string): Promise<void> => {
            const token = await getToken()
            return apiFetch<void>(`/urls/${encodeURIComponent(code)}`, token, { method: 'DELETE' })
        },

        stats: async (code: string): Promise<LinkStats> => {
            const token = await getToken()
            return apiFetch<LinkStats>(`/urls/${encodeURIComponent(code)}/stats`, token)
        },
        }),
        [getToken]
    )
}
