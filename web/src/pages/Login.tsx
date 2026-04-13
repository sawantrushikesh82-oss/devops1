import { useState } from 'react'
import api, { setAccessToken } from '../api/client'

export default function Login() {
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')

  const submit = async (e: React.FormEvent) => {
    e.preventDefault()
    const { data } = await api.post('/auth/login', { username, password })
    setAccessToken(data.access_token)
    window.location.href = '/users'
  }

  return <form onSubmit={submit} className="p-6 max-w-sm mx-auto space-y-2"><input className="border p-2 w-full" placeholder="Username" value={username} onChange={e=>setUsername(e.target.value)} /><input type="password" className="border p-2 w-full" placeholder="Password" value={password} onChange={e=>setPassword(e.target.value)} /><button className="bg-blue-600 text-white px-4 py-2">Login</button></form>
}
