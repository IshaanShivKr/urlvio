import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'

import { useLinksApi } from '../lib/links'
import type { CreateLinkInput } from '../types/link'

export const linkKeys = {
    all: ['links'] as const,
    lists: () => [...linkKeys.all, 'list'] as const,
    detail: (code: string) => [...linkKeys.all, 'detail', code] as const,
}

export function useLinksQuery() {
    const linksApi = useLinksApi()

    return useQuery({
        queryKey: linkKeys.lists(),
        queryFn: () => linksApi.list(),
        refetchInterval: 15_000,
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
