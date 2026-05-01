import { useMutation, useQuery } from '@tanstack/react-query'
import api from '../api/client'

export default function Sync(){
  const {data} = useQuery({queryKey:['sync-status'], queryFn:async()=> (await api.get('/sync/status')).data, refetchInterval:1000})
  const trigger = useMutation({mutationFn:()=> api.post('/sync/trigger')})
  const next = data?.next_run ? Math.max(0, Math.floor((new Date(data.next_run).getTime()-Date.now())/1000)) : 0
  return <div className="p-6"><h1>Sync</h1><p>Last Run: {String(data?.last_run||'never')}</p><p>Next run in: {next}s</p><button onClick={()=>trigger.mutate()}>Trigger Sync</button>{data?.running && <span className="ml-2 animate-spin">⏳</span>}</div>
}
