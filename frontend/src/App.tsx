import { Navigate, Route, Routes } from 'react-router-dom'

import { AppLayout } from './layouts/AppLayout'
import { Dashboard } from './pages/Dashboard'
import { Landing } from './pages/Landing'
import { SignIn } from './pages/SignIn'
import { SignUp } from './pages/SignUp'
import { UrlDetails } from './pages/UrlDetails'

function App() {
  return (
    <Routes>
      <Route path="/" element={<Landing />} />
      <Route path="/sign-in" element={<SignIn />} />
      <Route path="/sign-up" element={<SignUp />} />

      <Route element={<AppLayout />}>
        <Route path="/dashboard" element={<Dashboard />} />
        <Route path="/urls/:code" element={<UrlDetails />} />
      </Route>
      
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  )
}

export default App
