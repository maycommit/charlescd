import api from './api'

const path = '/workspaces'

export const getCircles = (workspaceId: string, clusterId: string) => api.get(`${path}/${workspaceId}/clusters/${clusterId}/circles`).then((res: any) => res.data)

export const getCircle = (workspaceId: string, clusterId: string, circleId: string) => api.get(`${path}/${workspaceId}/clusters/${clusterId}/circles/${circleId}`).then((res: any) => res.data)

export const deleteCircle = (workspaceId: string, clusterId: string, circleId: string) => api.delete(`${path}/${workspaceId}/clusters/${clusterId}/circles/${circleId}`).then((res: any) => res.data)

export const getManifest = (workspaceId: string, clusterId: string, name: string, version: string, group: string, kind: string) => api.get(`${path}/${workspaceId}/clusters/${clusterId}/manifests/${name}?kind=${kind}&group=${group}&version=${version}`).then((res: any) => res.data)

export const deploy = (workspaceId: string, clusterId: string, circleId: string, release: any) => api.post(`${path}/${workspaceId}/clusters/${clusterId}/circles/${circleId}/release`, release).then((res: any) => res.data)

export const undeploy = (workspaceId: string, clusterId: string, circleId: string) => api.delete(`${path}/${workspaceId}/clusters/${clusterId}/circles/${circleId}/release`).then((res: any) => res.data)

export const addProject = (workspaceId: string, clusterId: string, circleId: string, project: any) => api.post(`${path}/${workspaceId}/clusters/${clusterId}/circles/${circleId}/projects`, project).then((res: any) => res.data)

export const removeProject = (workspaceId: string, clusterId: string, circleId: string, projectName: string) => api.delete(`${path}/${workspaceId}/clusters/${clusterId}/circles/${circleId}/projects/${projectName}`).then((res: any) => res.data)

export const createCircle = (workspaceId: string, clusterId: string, data: any) =>
  api.post(`${path}/${workspaceId}/clusters/${clusterId}/circles`, data).then((res: any) => res.data)

export const getCircleTree = (workspaceId: string, clusterId: string, circleName: string) => api.get(`${path}/${workspaceId}/clusters/${clusterId}/circles/${circleName}/tree`).then(res => res.data)