import { useQuery } from '@tanstack/react-query'

import { useLinksApi } from '../lib/links'

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
    })
}
