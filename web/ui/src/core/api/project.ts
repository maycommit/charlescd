import api from './api'

const path = '/workspaces'

export const getProjects = (workspaceId: string, clusterId: string) => api.get(`${path}/${workspaceId}/clusters/${clusterId}/projects`).then((res: any) => res.data)

export const deleteProject = (workspaceId: string, clusterId: string, projectId: string) => api.delete(`${path}/${workspaceId}/clusters/${clusterId}/projects/${projectId}`).then((res: any) => res.data)

export const getProject = (name: string) => api.get(`/projects/${name}`).then(res => res.data)

export const createProject = (workspaceId: string, clusterId: string, data: any) =>
  api.post(`${path}/${workspaceId}/clusters/${clusterId}/projects`, data).then((res: any) => res.data)
