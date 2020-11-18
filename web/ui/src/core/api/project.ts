import api from './api'

export const getProjects = () => api.get('/projects').then(res => res.data)

export const createProject = (data: any) => api.post('/projects', data)