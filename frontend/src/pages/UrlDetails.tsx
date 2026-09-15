import { Clipboard } from 'lucide-react'
import { useEffect, useId, useRef, useState, type FormEvent } from 'react'
import { Link, useNavigate, useParams } from 'react-router-dom'
import { toast } from 'sonner'

import { ApiError } from '../lib/api'
import { validateUrl } from '../lib/validateUrl'
import {
    useDeleteLinkMutation,
    useLinkQuery,
    useLinkStatsQuery,
    useUpdateLinkMutation,
} from '../queries/links'

const focusRing =
    'focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent'

const dateFormatter = new Intl.DateTimeFormat('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
})

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

function errorMessageFor(error: unknown, fallback: string): string {
    return error instanceof ApiError ? error.message : fallback
}

export function UrlDetails() {
    const { code } = useParams<{ code: string }>()
    const navigate = useNavigate()

    const linkQuery = useLinkQuery(code)
    const statsQuery = useLinkStatsQuery(code)
    const updateLink = useUpdateLinkMutation()
    const deleteLink = useDeleteLinkMutation()

    const [isEditing, setIsEditing] = useState(false)
    const [destinationInput, setDestinationInput] = useState('')
    const [editError, setEditError] = useState<string | null>(null)
    const [isConfirmingDelete, setIsConfirmingDelete] = useState(false)

    const editErrorId = useId()
    const deleteDialogTitleId = useId()
    const deleteDialogRef = useRef<HTMLDialogElement>(null)

    useEffect(() => {
        const dialog = deleteDialogRef.current
        if (!dialog) return

        if (isConfirmingDelete && !dialog.open) {
        dialog.showModal()
        } else if (!isConfirmingDelete && dialog.open) {
        dialog.close()
        }
    }, [isConfirmingDelete])

    const backLink = (
        <Link to="/dashboard" className={`rounded text-sm text-muted transition-colors hover:text-ink ${focusRing}`}>
        Back to URLs
        </Link>
    )

    if (!code) {
        return (
        <div>
            {backLink}
            <div className="mt-6 rounded border border-zinc-200 bg-surface py-16 text-center">
            <p className="text-sm text-muted">This URL doesn&apos;t have a valid code.</p>
            </div>
        </div>
        )
    }

    const startEditing = () => {
        setDestinationInput(linkQuery.data?.url ?? '')
        setEditError(null)
        setIsEditing(true)
    }

    const cancelEditing = () => {
        if (updateLink.isPending) return
        setIsEditing(false)
        setEditError(null)
    }

    const handleSaveEdit = (event: FormEvent<HTMLFormElement>) => {
        event.preventDefault()
        if (updateLink.isPending) return

        const error = validateUrl(destinationInput)
        if (error) {
        setEditError(error)
        return
        }

        setEditError(null)

        updateLink.mutate(
        { code, input: { url: destinationInput.trim() } },
        {
            onSuccess: () => {
            toast.success('Destination updated')
            setIsEditing(false)
            },
            onError: (err) => {
            toast.error(errorMessageFor(err, 'Could not update the URL. Please try again.'))
            // Deliberately not exiting edit mode — the user's input stays put.
            },
        }
        )
    }

    const handleConfirmDelete = () => {
        if (deleteLink.isPending) return

        deleteLink.mutate(code, {
        onSuccess: () => {
            toast.success('URL deleted')
            navigate('/dashboard')
        },
        onError: (err) => {
            toast.error(errorMessageFor(err, 'Could not delete the URL. Please try again.'))
            setIsConfirmingDelete(false)
        },
        })
    }

    const isNotFound = linkQuery.error instanceof ApiError && linkQuery.error.status === 404

    return (
        <div>
        {backLink}
        <h1 className="mt-2 text-2xl font-semibold tracking-tight">URL Details</h1>

        {linkQuery.isPending && (
            <p className="py-16 text-center text-sm text-muted">Loading URL details…</p>
        )}

        {linkQuery.isError && isNotFound && (
            <div className="mt-6 rounded border border-zinc-200 bg-surface py-16 text-center">
            <p className="text-sm text-muted">This URL doesn&apos;t exist or has already been deleted.</p>
            <Link
                to="/dashboard"
                className={`mt-4 inline-block rounded border border-ink px-4 py-1.5 text-sm font-medium text-ink transition-colors hover:bg-ink hover:text-white ${focusRing}`}
            >
                Back to URLs
            </Link>
            </div>
        )}

        {linkQuery.isError && !isNotFound && (
            <div
            role="alert"
            className="mt-6 flex flex-col items-center gap-3 rounded border border-zinc-200 bg-surface py-16 text-center"
            >
            <p className="text-sm text-muted">
                {errorMessageFor(linkQuery.error, 'Something went wrong. Please try again.')}
            </p>
            <button
                type="button"
                onClick={() => void linkQuery.refetch()}
                className={`rounded border border-ink px-4 py-1.5 text-sm font-medium text-ink transition-colors hover:bg-ink hover:text-white ${focusRing}`}
            >
                Retry
            </button>
            </div>
        )}

        {linkQuery.isSuccess && (
            <div className="mt-6 max-w-2xl space-y-8">
            <section>
                <h2 className="text-sm font-medium text-muted">Short URL</h2>
                <div className="mt-1 flex items-center gap-2">
                <a
                    href={linkQuery.data.short_url}
                    target="_blank"
                    rel="noreferrer"
                    className={`rounded font-mono text-accent hover:underline ${focusRing}`}
                >
                    {displayShortUrl(linkQuery.data.short_url)}
                </a>
                <button
                    type="button"
                    onClick={() => handleCopy(linkQuery.data.short_url)}
                    aria-label={`Copy short URL ${displayShortUrl(linkQuery.data.short_url)}`}
                    className={`rounded p-1 text-muted transition-colors hover:text-ink ${focusRing}`}
                >
                    <Clipboard size={16} aria-hidden="true" />
                </button>
                </div>
            </section>

            <section>
                <h2 className="text-sm font-medium text-muted">Destination</h2>

                {!isEditing ? (
                <div className="mt-1 flex items-start justify-between gap-4">
                    <p className="min-w-0 flex-1 wrap-break-word text-sm text-ink">{linkQuery.data.url}</p>
                    <button
                    type="button"
                    onClick={startEditing}
                    className={`shrink-0 rounded text-sm text-muted transition-colors hover:text-ink ${focusRing}`}
                    >
                    Edit
                    </button>
                </div>
                ) : (
                <form onSubmit={handleSaveEdit} className="mt-2" noValidate>
                    <label htmlFor="destination-edit" className="sr-only">
                    Destination URL
                    </label>
                    <input
                    id="destination-edit"
                    type="url"
                    inputMode="url"
                    autoFocus
                    value={destinationInput}
                    onChange={(event) => {
                        setDestinationInput(event.target.value)
                        if (editError) setEditError(null)
                    }}
                    disabled={updateLink.isPending}
                    aria-invalid={editError ? true : undefined}
                    aria-describedby={editError ? editErrorId : undefined}
                    className={`w-full rounded border border-zinc-200 px-3 py-2 text-sm text-ink placeholder:text-muted disabled:bg-surface disabled:text-muted ${focusRing}`}
                    />
                    {editError && (
                    <p id={editErrorId} role="alert" className="mt-2 text-sm text-red-600">
                        {editError}
                    </p>
                    )}
                    <div className="mt-3 flex gap-3">
                    <button
                        type="submit"
                        disabled={updateLink.isPending}
                        className={`rounded border border-ink bg-ink px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-zinc-800 disabled:cursor-not-allowed disabled:opacity-60 ${focusRing}`}
                    >
                        {updateLink.isPending ? 'Saving…' : 'Save'}
                    </button>
                    <button
                        type="button"
                        onClick={cancelEditing}
                        disabled={updateLink.isPending}
                        className={`rounded border border-zinc-200 px-4 py-2 text-sm font-medium text-ink transition-colors hover:bg-surface disabled:cursor-not-allowed disabled:opacity-60 ${focusRing}`}
                    >
                        Cancel
                    </button>
                    </div>
                </form>
                )}
            </section>

            <section>
                <h2 className="text-sm font-medium text-muted">Created</h2>
                <p className="mt-1 text-sm text-ink">{formatDate(linkQuery.data.created_at)}</p>
            </section>

            <section>
                <h2 className="text-sm font-medium text-muted">Statistics</h2>

                {statsQuery.isPending && <p className="mt-1 text-sm text-muted">Loading statistics…</p>}

                {statsQuery.isError && (
                <p className="mt-1 text-sm text-muted">
                    {errorMessageFor(statsQuery.error, "Couldn't load statistics.")}
                </p>
                )}

                {statsQuery.isSuccess && (
                <dl className="mt-1 grid max-w-xs grid-cols-2 gap-4">
                    <div>
                    <dt className="text-xs text-muted">Total accesses</dt>
                    <dd className="mt-1 text-lg font-semibold tabular-nums">
                        {formatCount(statsQuery.data.access_count)}
                    </dd>
                    </div>
                    <div>
                    <dt className="text-xs text-muted">Last updated</dt>
                    <dd className="mt-1 text-sm text-ink">{formatDate(statsQuery.data.updated_at)}</dd>
                    </div>
                </dl>
                )}
            </section>

            <section className="border-t border-zinc-200 pt-6">
                <h2 className="text-sm font-medium text-muted">Danger zone</h2>
                <p className="mt-1 text-sm text-muted">Deleting this URL cannot be undone.</p>
                <button
                type="button"
                onClick={() => setIsConfirmingDelete(true)}
                className={`mt-3 rounded border border-zinc-200 px-4 py-2 text-sm font-medium text-red-600 transition-colors hover:bg-red-50 ${focusRing}`}
                >
                Delete URL
                </button>
            </section>
            </div>
        )}

        <dialog
            ref={deleteDialogRef}
            aria-labelledby={deleteDialogTitleId}
            aria-modal="true"
            className="m-auto w-[calc(100vw-2rem)] max-w-md rounded border border-zinc-200 bg-white p-0 text-ink backdrop:bg-zinc-950/40"
            onClose={() => setIsConfirmingDelete(false)}
            onCancel={(event) => {
            if (deleteLink.isPending) event.preventDefault()
            }}
            onClick={(event) => {
            if (event.target === deleteDialogRef.current && !deleteLink.isPending) {
                setIsConfirmingDelete(false)
            }
            }}
        >
            <div className="p-6">
            <h2 id={deleteDialogTitleId} className="text-lg font-semibold tracking-tight">
                Delete this URL?
            </h2>
            <p className="mt-2 text-sm text-muted">
                This action cannot be undone. The short URL will stop working immediately.
            </p>
            <div className="mt-6 flex justify-end gap-3">
                <button
                type="button"
                onClick={() => setIsConfirmingDelete(false)}
                disabled={deleteLink.isPending}
                className={`rounded border border-zinc-200 px-4 py-2 text-sm font-medium text-ink transition-colors hover:bg-surface disabled:cursor-not-allowed disabled:opacity-60 ${focusRing}`}
                >
                Cancel
                </button>
                <button
                type="button"
                onClick={handleConfirmDelete}
                disabled={deleteLink.isPending}
                className={`rounded border border-red-600 bg-red-600 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-red-700 disabled:cursor-not-allowed disabled:opacity-60 ${focusRing}`}
                >
                {deleteLink.isPending ? 'Deleting…' : 'Delete'}
                </button>
            </div>
            </div>
        </dialog>
        </div>
    )
}
