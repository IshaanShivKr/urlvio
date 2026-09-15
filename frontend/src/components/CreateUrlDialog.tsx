import { useEffect, useId, useRef, useState, type FormEvent } from 'react'
import { toast } from 'sonner'

import { ApiError } from '../lib/api'
import { validateUrl } from '../lib/validateUrl'
import { useCreateLinkMutation } from '../queries/links'

const focusRing =
    'focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent'

interface CreateUrlDialogProps {
    open: boolean
    onOpenChange: (open: boolean) => void
}

export function CreateUrlDialog({ open, onOpenChange }: CreateUrlDialogProps) {
    const dialogRef = useRef<HTMLDialogElement>(null)
    const [url, setUrl] = useState('')
    const [validationError, setValidationError] = useState<string | null>(null)
    const titleId = useId()
    const errorId = useId()

    const createLink = useCreateLinkMutation()

    useEffect(() => {
        const dialog = dialogRef.current
        if (!dialog) return

        if (open && !dialog.open) {
        dialog.showModal()
        } else if (!open && dialog.open) {
        dialog.close()
        }
    }, [open])

    function resetAndClose() {
        setUrl('')
        setValidationError(null)
        createLink.reset()
        onOpenChange(false)
    }

    function handleSubmit(event: FormEvent<HTMLFormElement>) {
        event.preventDefault()

        if (createLink.isPending) return

        const error = validateUrl(url)
        if (error) {
        setValidationError(error)
        return
        }

        setValidationError(null)

        createLink.mutate(
        { url: url.trim() },
        {
            onSuccess: (link) => {
            toast.success('Short URL created', { description: link.short_url })
            resetAndClose()
            },
            onError: (err) => {
            const message =
                err instanceof ApiError ? err.message : 'Could not create the URL. Please try again.'
            toast.error(message)
            },
        }
        )
    }

    return (
        <dialog
        ref={dialogRef}
        aria-labelledby={titleId}
        aria-modal="true"
        className="m-auto w-[calc(100vw-2rem)] max-w-md rounded border border-zinc-200 bg-white p-0 text-ink backdrop:bg-zinc-950/40"
        onClose={resetAndClose}
        onCancel={(event) => {
            if (createLink.isPending) {
            event.preventDefault()
            }
        }}
        onClick={(event) => {
            if (event.target === dialogRef.current && !createLink.isPending) {
            resetAndClose()
            }
        }}
        >
        <form onSubmit={handleSubmit} className="p-6" noValidate>
            <div className="flex items-start justify-between gap-4">
            <h2 id={titleId} className="text-lg font-semibold tracking-tight">
                Create URL
            </h2>
            <button
                type="button"
                onClick={resetAndClose}
                disabled={createLink.isPending}
                aria-label="Close"
                className={`rounded text-muted transition-colors hover:text-ink disabled:cursor-not-allowed disabled:opacity-60 ${focusRing}`}
            >
                ×
            </button>
            </div>

            <div className="mt-5">
            <label htmlFor="destination-url" className="block text-sm font-medium text-ink">
                Destination URL
            </label>
            <input
                id="destination-url"
                type="url"
                inputMode="url"
                autoFocus
                required
                value={url}
                onChange={(event) => {
                setUrl(event.target.value)
                if (validationError) setValidationError(null)
                }}
                placeholder="https://example.com/your/long/path"
                aria-invalid={validationError ? true : undefined}
                aria-describedby={validationError ? errorId : undefined}
                disabled={createLink.isPending}
                className={`mt-2 w-full rounded border border-zinc-200 px-3 py-2 text-sm text-ink placeholder:text-muted disabled:bg-surface disabled:text-muted ${focusRing}`}
            />
            {validationError && (
                <p id={errorId} role="alert" className="mt-2 text-sm text-red-600">
                {validationError}
                </p>
            )}
            </div>

            <div className="mt-6 flex justify-end gap-3">
            <button
                type="button"
                onClick={resetAndClose}
                disabled={createLink.isPending}
                className={`rounded border border-zinc-200 px-4 py-2 text-sm font-medium text-ink transition-colors hover:bg-surface disabled:cursor-not-allowed disabled:opacity-60 ${focusRing}`}
            >
                Cancel
            </button>
            <button
                type="submit"
                disabled={createLink.isPending}
                className={`rounded border border-ink bg-ink px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-zinc-800 disabled:cursor-not-allowed disabled:opacity-60 ${focusRing}`}
            >
                {createLink.isPending ? 'Creating…' : 'Create'}
            </button>
            </div>
        </form>
        </dialog>
    )
}
