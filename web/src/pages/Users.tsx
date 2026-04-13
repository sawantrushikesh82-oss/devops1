import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useState } from 'react'
import api from '../api/client'

export default function Users() {
  const qc = useQueryClient()
  const [newUser, setNewUser] = useState({ username:'', passwordHash:'', userType:'external', role:'viewer', publicKey:'' })
  const { data=[] } = useQuery({ queryKey:['users'], queryFn: async()=> (await api.get('/users')).data })
  const createM = useMutation({ mutationFn:()=> api.post('/users', newUser), onSuccess:()=> qc.invalidateQueries({queryKey:['users']}) })
  const delM = useMutation({ mutationFn:(id:string)=> api.delete(`/users/${id}`), onSuccess:()=> qc.invalidateQueries({queryKey:['users']}) })
  return <div className="p-6"><h1 className="text-xl mb-4">Users</h1><table className="w-full border"><thead><tr><th>Username</th><th>Type</th><th>Role</th><th/></tr></thead><tbody>{data.map((u:any)=><tr key={u.id}><td>{u.username}</td><td><span className="px-2 py-1 bg-slate-200 rounded">{u.user_type}</span></td><td>{u.role}</td><td><button onClick={()=>delM.mutate(u.id)}>Delete</button></td></tr>)}</tbody></table><div className="mt-4 space-x-2"><input className="border p-1" placeholder="username" onChange={e=>setNewUser({...newUser,username:e.target.value})}/><input className="border p-1" placeholder="password" onChange={e=>setNewUser({...newUser,passwordHash:e.target.value})}/><select onChange={e=>setNewUser({...newUser,userType:e.target.value})}><option>external</option><option>internal</option></select><button className="bg-green-600 text-white px-2" onClick={()=>createM.mutate()}>Add User</button></div></div>
}
