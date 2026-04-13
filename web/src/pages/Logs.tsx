import { useQuery } from '@tanstack/react-query'
import { useState } from 'react'
import api from '../api/client'

export default function Logs(){
  const [page,setPage]=useState(1); const [status,setStatus]=useState(''); const [username,setUsername]=useState('')
  const {data=[]}=useQuery({queryKey:['logs',page,status,username],queryFn:async()=> (await api.get('/logs',{params:{page,limit:20,status,username}})).data, refetchInterval:30000})
  return <div className="p-6"><h1>Logs</h1><div><input placeholder="username" onChange={e=>setUsername(e.target.value)}/><select onChange={e=>setStatus(e.target.value)}><option value="">all</option><option>success</option><option>failed</option><option>access_denied</option></select></div><table><tbody>{data.map((l:any)=><tr key={l.id}><td>{l.created_at}</td><td>{l.username}</td><td>{l.filename}</td><td>{l.status}</td></tr>)}</tbody></table><button onClick={()=>setPage(Math.max(1,page-1))}>Prev</button><button onClick={()=>setPage(page+1)}>Next</button></div>
}
