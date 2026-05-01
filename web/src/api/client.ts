import axios from 'axios'

let accessToken = ''
const api = axios.create({ baseURL: '/api', withCredentials: true })

export const setAccessToken = (t: string) => { accessToken = t }

api.interceptors.request.use((cfg) => {
  if (accessToken) cfg.headers.Authorization = `Bearer ${accessToken}`
  return cfg
})

api.interceptors.response.use((r) => r, async (error) => {
  const original = error.config
  if (error.response?.status === 401 && !original._retry) {
    original._retry = true
    const { data } = await api.post('/auth/refresh')
    setAccessToken(data.access_token)
    original.headers.Authorization = `Bearer ${data.access_token}`
    return api(original)
  }
  return Promise.reject(error)
})

export default api
