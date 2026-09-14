import { Link, Outlet } from 'react-router-dom'

export function AppLayout() {
    return (
        <div className="min-h-screen bg-white text-zinc-950">
        <header className="border-b border-zinc-200">
            <div className="mx-auto flex h-14 max-w-6xl items-center justify-between px-6">
            <Link to="/dashboard" className="font-semibold tracking-tight">
                URLVio
            </Link>

            <nav className="flex items-center gap-5 text-sm text-zinc-600">
                <Link to="/dashboard" className="hover:text-zinc-950">
                Dashboard
                </Link>
            </nav>
            </div>
        </header>

        <main className="mx-auto max-w-6xl px-6 py-8">
            <Outlet />
        </main>
        </div>
    )
}
