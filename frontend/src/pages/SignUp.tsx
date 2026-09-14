import { SignUp as ClerkSignUp } from '@clerk/clerk-react'

export function SignUp() {
    return (
        <div className="flex min-h-screen items-center justify-center">
        <ClerkSignUp routing="path" path="/sign-up" signInUrl="/sign-in" forceRedirectUrl="/dashboard" />
        </div>
    )
}
