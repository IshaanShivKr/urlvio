import { useMemo, useState } from 'react'
import { Link } from 'react-router-dom'
import { toast } from 'sonner'
import { Clipboard } from 'lucide-react'

import { CreateUrlDialog } from '../components/CreateUrlDialog'
import { ApiError } from '../lib/api'
import { useLinksQuery } from '../queries/links'

const dateFormatter = new Intl.DateTimeFormat('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
})

const focusRing =
    'focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent'

function formatDate(iso: string): string {
    return dateFormatter.format(new Date(iso))
}

function formatCount(count: number): string {
    return count.toLocaleString('en-US')
}

function displayShortUrl(shortUrl: string): string {
    return shortUrl.replace(/^https?:\/\//, '')
}

function handleCopy(shortUrl: string) {
    navigator.clipboard
        .writeText(shortUrl)
        .then(() => toast.success('Copied to clipboard'))
        .catch(() => toast.error('Could not copy to clipboard'))
}

export function Dashboard() {
    const query = useLinksQuery()
    const [search, setSearch] = useState('')
    const [isCreateOpen, setIsCreateOpen] = useState(false)

    const filteredLinks = useMemo(() => {
        const links = query.data ?? []

        const trimmed = search.trim().toLowerCase()
        if (!trimmed) return links

        return links.filter(
        (link) => link.code.toLowerCase().includes(trimmed) || link.url.toLowerCase().includes(trimmed)
        )
    }, [query.data, search])

    const hasLinks = (query.data?.length ?? 0) > 0
    const showNoResults = hasLinks && filteredLinks.length === 0
    const showResults = filteredLinks.length > 0

    const errorMessage =
        query.error instanceof ApiError ? query.error.message : 'Something went wrong. Please try again.'

    return (
        <div>
        <div className="flex flex-wrap items-center justify-between gap-4">
            <h1 className="text-2xl font-semibold tracking-tight">Your URLs</h1>
            <button
            type="button"
            onClick={() => setIsCreateOpen(true)}
            className={`rounded border border-ink bg-ink px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-zinc-800 ${focusRing}`}
            >
            Create URL
            </button>
        </div>

        <CreateUrlDialog open={isCreateOpen} onOpenChange={setIsCreateOpen} />

        <div className="mt-6">
            <label htmlFor="url-search" className="sr-only">
            Search your URLs
            </label>
            <input
            id="url-search"
            type="search"
            value={search}
            onChange={(event) => setSearch(event.target.value)}
            placeholder="Search by code or destination…"
            disabled={query.isPending}
            className={`w-full max-w-sm rounded border border-zinc-200 px-3 py-2 text-sm text-ink placeholder:text-muted disabled:bg-surface disabled:text-muted ${focusRing}`}
            />
        </div>

        <div className="mt-6" aria-live="polite">
            {query.isPending && (
            <p className="py-16 text-center text-sm text-muted">Loading your URLs…</p>
            )}

            {query.isError && (
            <div
                role="alert"
                className="flex flex-col items-center gap-3 rounded border border-zinc-200 bg-surface py-16 text-center"
            >
                <p className="text-sm text-muted">{errorMessage}</p>
                <button
                type="button"
                onClick={() => void query.refetch()}
                className={`rounded border border-ink px-4 py-1.5 text-sm font-medium text-ink transition-colors hover:bg-ink hover:text-white ${focusRing}`}
                >
                Retry
                </button>
            </div>
            )}

            {query.isSuccess && !hasLinks && (
            <div className="rounded border border-zinc-200 bg-surface py-16 text-center">
                <p className="text-sm text-muted">You haven&apos;t created any URLs yet.</p>
            </div>
            )}

            {showNoResults && (
            <div className="rounded border border-zinc-200 bg-surface py-16 text-center">
                <p className="text-sm text-muted">No URLs match &ldquo;{search}&rdquo;.</p>
            </div>
            )}

            {showResults && (
            <>
                {/* Table for sm and up */}
                <table className="hidden w-full border-collapse text-sm sm:table">
                <thead>
                    <tr className="border-b border-zinc-200 text-left text-xs text-muted">
                    <th scope="col" className="py-2 pr-4 font-medium">
                        Short URL
                    </th>
                    <th scope="col" className="py-2 pr-4 font-medium">
                        Destination
                    </th>
                    <th scope="col" className="py-2 pr-4 font-medium">
                        Accesses
                    </th>
                    <th scope="col" className="py-2 pr-4 font-medium">
                        Created
                    </th>
                    <th scope="col" className="py-2 font-medium">
                        <span className="sr-only">Actions</span>
                    </th>
                    </tr>
                </thead>
                <tbody>
                    {filteredLinks.map((link) => (
                    <tr key={link.code} className="border-b border-zinc-100">
                        <td className="py-3 pr-4 font-mono">
                        <div className="flex items-center gap-2">
                            <a
                            href={link.short_url}
                            target="_blank"
                            rel="noreferrer"
                            className={`rounded text-accent hover:underline ${focusRing}`}
                            >
                            {displayShortUrl(link.short_url)}
                            </a>
                            <button
                                type="button"
                                onClick={() => handleCopy(link.short_url)}
                                aria-label={`Copy short URL ${displayShortUrl(link.short_url)}`}
                                className={`rounded text-muted hover:text-ink ${focusRing}`}
                            >
                                <Clipboard size={16} />
                            </button>
                        </div>
                        </td>
                        <td className="max-w-xs truncate py-3 pr-4 text-muted" title={link.url}>
                        {link.url}
                        </td>
                        <td className="py-3 pr-4 tabular-nums">{formatCount(link.access_count)}</td>
                        <td className="py-3 pr-4 text-muted">{formatDate(link.created_at)}</td>
                        <td className="py-3 text-right">
                        <Link to={`/urls/${link.code}`} className={`rounded text-muted hover:text-ink ${focusRing}`}>
                            View
                        </Link>
                        </td>
                    </tr>
                    ))}
                </tbody>
                </table>

                {/* Stacked list below sm */}
                <ul className="divide-y divide-zinc-100 sm:hidden">
                {filteredLinks.map((link) => (
                    <li key={link.code} className="py-4">
                    <div className="flex items-center justify-between gap-3">
                        <div className="flex items-center gap-2">
                        <a
                            href={link.short_url}
                            target="_blank"
                            rel="noreferrer"
                            className={`rounded font-mono text-accent hover:underline ${focusRing}`}
                        >
                            {displayShortUrl(link.short_url)}
                        </a>
                        <button
                            type="button"
                            onClick={() => handleCopy(link.short_url)}
                            aria-label={`Copy short URL ${displayShortUrl(link.short_url)}`}
                            className={`rounded text-muted hover:text-ink ${focusRing}`}
                        >
                            <Clipboard size={16} />
                        </button>
                        </div>
                        <Link to={`/urls/${link.code}`} className={`rounded text-sm text-muted hover:text-ink ${focusRing}`}>
                        View
                        </Link>
                    </div>
                    <p className="mt-1 truncate text-sm text-muted" title={link.url}>
                        {link.url}
                    </p>
                    <div className="mt-2 flex items-center gap-4 text-xs text-muted">
                        <span>{formatCount(link.access_count)} accesses</span>
                        <span>{formatDate(link.created_at)}</span>
                    </div>
                    </li>
                ))}
                </ul>
            </>
            )}
        </div>
        </div>
    )
}
