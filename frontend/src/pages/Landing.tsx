import { Link } from 'react-router-dom'

const API_DOCS_URL = 'https://urlvio.onrender.com/docs/api'
const GITHUB_URL = 'https://github.com/IshaanShivKr/urlvio'

const focusRing =
    'focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent'

export function Landing() {
    return (
        <div className="flex min-h-screen flex-col bg-white text-ink">
        <header className="border-b border-zinc-200">
            <div className="mx-auto flex h-16 max-w-5xl items-center justify-between px-6">
            <Link to="/" className={`rounded font-mono text-base font-semibold tracking-tight ${focusRing}`}>
                URLVio
            </Link>

            <nav className="flex items-center gap-6 text-sm">
                <a
                href={API_DOCS_URL}
                target="_blank"
                rel="noreferrer"
                className={`hidden rounded text-muted transition-colors hover:text-ink sm:inline ${focusRing}`}
                >
                API Docs
                </a>
                <Link to="/sign-in" className={`rounded text-muted transition-colors hover:text-ink ${focusRing}`}>
                Sign in
                </Link>
                <Link
                to="/sign-up"
                className={`rounded border border-ink bg-ink px-3.5 py-1.5 font-medium text-white transition-colors hover:bg-zinc-800 ${focusRing}`}
                >
                Get started
                </Link>
            </nav>
            </div>
        </header>

        <main className="flex-1">
            <section className="mx-auto grid max-w-5xl gap-14 px-6 py-20 sm:py-28 lg:grid-cols-[1.1fr_1fr] lg:items-center">
            <div>
                <h1 className="text-4xl font-semibold leading-tight tracking-tight sm:text-5xl">
                Shorten URLs. Track every redirect.
                </h1>
                <p className="mt-5 max-w-md text-base leading-relaxed text-muted">
                Create short, permanent links through a REST API or the dashboard. Every
                link keeps its own access count — no analytics suite required.
                </p>
                <div className="mt-8 flex flex-wrap items-center gap-x-6 gap-y-3">
                <Link
                    to="/sign-up"
                    className={`rounded border border-ink bg-ink px-5 py-2.5 text-sm font-medium text-white transition-colors hover:bg-zinc-800 ${focusRing}`}
                >
                    Get started
                </Link>
                <a
                    href={API_DOCS_URL}
                    target="_blank"
                    rel="noreferrer"
                    className={`rounded text-sm font-medium text-muted transition-colors hover:text-ink ${focusRing}`}
                >
                    Read the API docs
                </a>
                </div>
            </div>

            <div>
                <div className="rounded border border-zinc-200 bg-surface p-5 font-mono text-sm">
                <p className="truncate text-muted">
                    https://example.com/blog/2024/09/why-short-links-still-matter
                </p>
                <div className="my-3 flex items-center gap-3 text-xs text-muted" aria-hidden="true">
                    <span className="h-px flex-1 bg-zinc-200" />
                    <span>302</span>
                    <span className="h-px flex-1 bg-zinc-200" />
                </div>
                <p className="text-ink">
                    urlvio.onrender.com/<span className="text-accent">aB3x9K</span>
                </p>
                </div>
                <p className="mt-3 text-xs text-muted">
                Every redirect increments the link's access count.
                </p>
            </div>
            </section>
        </main>

        <footer className="border-t border-zinc-200">
            <div className="mx-auto flex max-w-5xl flex-col gap-3 px-6 py-8 text-sm text-muted sm:flex-row sm:items-center sm:justify-between">
            <span className="font-mono text-ink">URLVio</span>
            <a
                href={GITHUB_URL}
                target="_blank"
                rel="noreferrer"
                className={`rounded transition-colors hover:text-ink ${focusRing}`}
            >
                GitHub
            </a>
            </div>
        </footer>
        </div>
    )
}
