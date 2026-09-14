export interface Link {
    code: string
    short_url: string
    url: string
    created_at: string
    updated_at: string
}

export interface LinkListItem extends Link {
    access_count: number
}

export interface LinkStats {
    code: string
    url: string
    access_count: number
    created_at: string
    updated_at: string
}

export interface CreateLinkInput {
    url: string
}

export interface UpdateLinkInput {
    url: string
}