import api from './api'

export const getCircles = () => api.get('/circles').then(res => res.data)

export const createCircle = (data: any) => api.post('/circles', data)

export const deploy = (name: string, data: any) => api.post(`/circles/${name}/deploy`, data)