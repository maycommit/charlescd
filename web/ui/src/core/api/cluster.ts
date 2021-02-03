import api from './api'

const path = '/workspaces'

export const clusterHealth = (workspaceId: string, clusterId: string) =>
  api.get(`${path}/${workspaceId}/clusters/${clusterId}/health`).then((res: any) => res.data)

export const createCluster = (workspaceId: string, data: any) =>
  api.post(`${path}/${workspaceId}/clusters`, data).then((res: any) => res.data)

export const deleteCluster = (workspaceId: string, clusterId: string) => api.delete(`${path}/${workspaceId}/clusters/${clusterId}`).then((res: any) => res.data)
