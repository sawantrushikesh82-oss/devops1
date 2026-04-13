import { useQuery } from '@tanstack/react-query'
import api from '../api/client'

export default function Health(){
  const {data} = useQuery({queryKey:['health'], queryFn:async()=> (await api.get('/health')).data, refetchInterval:5000})
  const card=(name:string,val:any)=><div className="border p-3 rounded"><h3>{name}</h3><p className={val?'text-green-600':'text-red-600'}>{String(val)}</p></div>
  return <div className="p-6 grid grid-cols-4 gap-4">{card('API',true)}{card('DB',data?.db)}{card('ext-SFTP',data?.ext_sftp)}{card('int-SFTP',data?.int_sftp)}</div>
}
