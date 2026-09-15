import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'

import { useLinksApi } from '../lib/links'
import type { CreateLinkInput, Link, UpdateLinkInput } from '../types/link'

export const linkKeys = {
    all: ['links'] as const,
    lists: () => [...linkKeys.all, 'list'] as const,
    detail: (code: string) => [...linkKeys.all, 'detail', code] as const,
    stats: (code: string) => [...linkKeys.all, 'stats', code] as const,
}

export function useLinksQuery() {
    const linksApi = useLinksApi()

    return useQuery({
        queryKey: linkKeys.lists(),
        queryFn: () => linksApi.list(),
        refetchInterval: 15_000,
    })
}

export function useLinkQuery(code: string | undefined) {
    const linksApi = useLinksApi()

    return useQuery({
        queryKey: linkKeys.detail(code ?? ''),
        enabled: Boolean(code),
        queryFn: () => {
        if (!code) {
            throw new Error('useLinkQuery called without a code')
        }
        return linksApi.get(code)
        },
    })
}

export function useLinkStatsQuery(code: string | undefined) {
    const linksApi = useLinksApi()

    return useQuery({
        queryKey: linkKeys.stats(code ?? ''),
        enabled: Boolean(code),
        queryFn: () => {
        if (!code) {
            throw new Error('useLinkStatsQuery called without a code')
        }
        return linksApi.stats(code)
        },
    })
}

export function useCreateLinkMutation() {
    const linksApi = useLinksApi()
    const queryClient = useQueryClient()

    return useMutation({
        mutationFn: (input: CreateLinkInput) => linksApi.create(input),
        onSuccess: () => {
        void queryClient.invalidateQueries({ queryKey: linkKeys.lists() })
        },
    })
}

export function useUpdateLinkMutation() {
    const linksApi = useLinksApi()
    const queryClient = useQueryClient()

    return useMutation({
        mutationFn: ({ code, input }: { code: string; input: UpdateLinkInput }) => linksApi.update(code, input),
        onSuccess: (link: Link, variables) => {
        queryClient.setQueryData(linkKeys.detail(variables.code), link)
        void queryClient.invalidateQueries({ queryKey: linkKeys.lists() })
        },
    })
}

export function useDeleteLinkMutation() {
    const linksApi = useLinksApi()
    const queryClient = useQueryClient()

    return useMutation({
        mutationFn: (code: string) => linksApi.remove(code),
        onSuccess: (_data, code) => {
        queryClient.removeQueries({ queryKey: linkKeys.detail(code) })
        queryClient.removeQueries({ queryKey: linkKeys.stats(code) })
        void queryClient.invalidateQueries({ queryKey: linkKeys.lists() })
        },
    })
}
