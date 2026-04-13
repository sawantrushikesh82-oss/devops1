import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { useState } from 'react'
import api from '../api/client'

export default function Mappings(){
  const qc=useQueryClient(); const [m,setM]=useState({external_user:'',internal_user:'',flow_type:'send_receive'})
  const {data=[]}=useQuery({queryKey:['mappings'],queryFn:async()=> (await api.get('/mappings')).data})
  const add=useMutation({mutationFn:()=>api.post('/mappings',m),onSuccess:()=>qc.invalidateQueries({queryKey:['mappings']})})
  const del=useMutation({mutationFn:(id:string)=>api.delete(`/mappings/${id}`),onSuccess:()=>qc.invalidateQueries({queryKey:['mappings']})})
  return <div className="p-6"><h1 className="text-xl">Mappings</h1><table><tbody>{data.map((x:any)=><tr key={x.id}><td>{x.external_user}</td><td>↔</td><td>{x.internal_user}</td><td>{x.flow_type}</td><td><button onClick={()=>del.mutate(x.id)}>Delete</button></td></tr>)}</tbody></table><div className="space-x-2"><input placeholder="external" onChange={e=>setM({...m,external_user:e.target.value})}/><input placeholder="internal" onChange={e=>setM({...m,internal_user:e.target.value})}/><select onChange={e=>setM({...m,flow_type:e.target.value})}><option>send_receive</option><option>send_only</option><option>receive_only</option></select><button onClick={()=>add.mutate()}>Add</button></div></div>
}
