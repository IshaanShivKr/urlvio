import { ClerkProvider } from '@clerk/clerk-react'
import { QueryClientProvider } from '@tanstack/react-query'
import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { BrowserRouter } from 'react-router-dom'

import App from './App.tsx'
import './index.css'
import { queryClient } from './lib/queryClient'
import { Toaster } from 'sonner'

const clerkPublishableKey = import.meta.env.VITE_CLERK_PUBLISHABLE_KEY

if (!clerkPublishableKey) {
  throw new Error('Missing VITE_CLERK_PUBLISHABLE_KEY environment variable')
}

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <BrowserRouter>
      <QueryClientProvider client={queryClient}>
        <ClerkProvider publishableKey={clerkPublishableKey} afterSignOutUrl="/">
          <Toaster />
          <App />
        </ClerkProvider>
      </QueryClientProvider>
    </BrowserRouter>
  </StrictMode>
)
