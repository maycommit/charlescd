import api from './api'

const path = '/workspaces'

export const getWorkspaces = () => api.get(path).then((res: any) => res.data)

export const createWorkspace = (data: any) => api.post(path, data).then((res: any) => res.data)

export const deleteWorkspace = (workspaceId: string) => api.delete(`${path}/${workspaceId}`).then((res: any) => res.data)
