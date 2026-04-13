import { BrowserRouter, Navigate, Route, Routes } from 'react-router-dom'
import Login from './pages/Login'
import Users from './pages/Users'
import Mappings from './pages/Mappings'
import Logs from './pages/Logs'
import Sync from './pages/Sync'
import Health from './pages/Health'

export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Navigate to="/login" />} />
        <Route path="/login" element={<Login />} />
        <Route path="/users" element={<Users />} />
        <Route path="/mappings" element={<Mappings />} />
        <Route path="/logs" element={<Logs />} />
        <Route path="/sync" element={<Sync />} />
        <Route path="/health" element={<Health />} />
      </Routes>
    </BrowserRouter>
  )
}
