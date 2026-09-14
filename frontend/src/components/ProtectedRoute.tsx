import { useAuth } from '@clerk/clerk-react'
import { Navigate, Outlet, useLocation } from 'react-router-dom'

export function ProtectedRoute() {
    const { isLoaded, isSignedIn } = useAuth()
    const location = useLocation()

    if (!isLoaded) {
        return (
        <div className="flex min-h-screen items-center justify-center text-sm text-zinc-500">
            Loading…
        </div>
        )
    }

    if (!isSignedIn) {
        return <Navigate to="/sign-in" replace state={{ from: location }} />
    }

    return <Outlet />
}
