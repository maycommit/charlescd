import React, { useEffect, useLayoutEffect, useRef, useState } from 'react'
import dagreD3 from 'dagre-d3'
import * as d3 from 'd3'
import { Route, useHistory, useLocation, useParams, useRouteMatch } from 'react-router-dom'
import { deploy, getCircle, getCircleTree, undeploy, removeProject, addProject } from '../../core/api/circle'
import './style.scss'
import { getTree } from './helper'
import NodeSidebar from './NodeSidebar'
import Sidebar from './Sidebar'
import { ROUTES_PREFIX } from '../../core/constants/routes'
import { Modal, ModalBody, ModalHeader } from 'reactstrap'
import CircleReleaseForm from './CircleRelease'
import CircleProjectForm from './CircleProject'

const ModalDeploy = () => {
  const history = useHistory()
  const match = useRouteMatch()
  const { workspaceId, clusterId, circleId } = useParams<any>()

  const toggle = () => history.goBack()

  const handleSubmit = (data: any) => {
    deploy(workspaceId, clusterId, circleId, data)
      .then(() => {
        console.log('Circle CREATE SUCCESS')
      })

    toggle()
  }

  return (
    <Modal isOpen={true} toggle={toggle} className="" size="lg">
      <ModalHeader toggle={toggle}>Deploy in circle</ModalHeader>
      <ModalBody>
        <CircleReleaseForm onSubmit={handleSubmit} />
      </ModalBody>
    </Modal>
  )
}

const ModalProject = () => {
  const history = useHistory()
  const match = useRouteMatch()
  const { workspaceId, clusterId, circleId } = useParams<any>()
  const [projects, setProjects] = useState([])

  const toggle = () => history.goBack()

  useEffect(() => {
    getCircle(workspaceId, clusterId, circleId)
      .then((data: any) => setProjects(data?.release?.projects))
  }, [])

  const handleSubmit = (data: any) => {
    addProject(workspaceId, clusterId, circleId, data)
      .then(() => {
        console.log('Circle CREATE SUCCESS')
      })

    toggle()
  }

  return (
    <Modal isOpen={true} toggle={toggle} className="" size="lg">
      <ModalHeader toggle={toggle}>Add project to release</ModalHeader>
      <ModalBody>
        <CircleProjectForm selectedProjects={projects} onSubmit={handleSubmit} />
      </ModalBody>
    </Modal>
  )
}

const Diagram = () => {
  const [diagramData, setDiagramData] = useState<any>({ nodes: [], edges: [] })

  const [selectedNode, setSelectedNode] = useState<any>()
  const { workspaceId, clusterId, circleId } = useParams<any>()
  const svg = useRef<any>()
  const innerG = useRef<any>()

  const requestCircleTree = () => getCircleTree(workspaceId, clusterId, circleId)
    .then((data: any) => {
      const [currentNodes, currentEdges] = getTree(data, circleId)
      setDiagramData({ nodes: currentNodes, edges: currentEdges })
    })

  useLayoutEffect(() => {
    requestCircleTree()
    const interval = setInterval(() => requestCircleTree(), 3000)

    return () => clearInterval(interval)
  }, [])


  useEffect(() => {
    let g = new dagreD3.graphlib.Graph().setGraph({
      marginx: 350,
      marginy: 15,
      rankdir: "LR",
      ranksep: 55,
      nodesep: 15
    })
    let render = new dagreD3.render()
    let innerSvg: any = d3.select(svg.current)
    let inner: any = d3.select(innerG.current)
    let zoom = d3.zoom().on('zoom', () => inner.attr('transform', (d3 as any)?.event.transform))
    innerSvg.call(zoom)

    diagramData?.nodes
      ?.sort((a: any, b: any) => {
        return a?.meta?.name?.localeCompare(b?.meta?.name) && a?.meta?.kind?.localeCompare(b?.meta?.kind)
      })
      ?.map((node: any) => {
        g.setNode(node.name, { ...node })
      })

    diagramData?.edges
      ?.sort((a: any, b: any) => {
        return a.source.localeCompare(b.source) && a.target.localeCompare(b.target)
      })
      ?.map((edge: any) => {
        g.setEdge(edge.source, edge.target, { ...edge })
      })

    render(inner, (g as any))

    innerSvg.selectAll('g.node').on('click', (id: any) => {
      setSelectedNode({ id, ...g.node(id) })
    })
  }, [diagramData])


  return (
    <>
      <svg ref={svg}><g ref={innerG} /></svg>
      {selectedNode && <NodeSidebar node={selectedNode} onClose={() => setSelectedNode(null)} />}
    </>
  )
}

const Tree = () => {
  const location = useLocation()
  const [circle, setCircle] = useState(null)
  const { workspaceId, clusterId, circleId } = useParams<any>()

  const requestCircle = () => getCircle(workspaceId, clusterId, circleId)
    .then(data => setCircle(data))

  useEffect(() => {
    requestCircle()
  }, [location])

  const handleUndeploy = () => undeploy(workspaceId, clusterId, circleId)
    .then(() => requestCircle())

  const handleRemoveProject = (projectName: string) => removeProject(workspaceId, clusterId, circleId, projectName)
    .then(() => requestCircle())

  return (
    <div className="tree">
      {circle && <Sidebar circle={circle} onUndeploy={handleUndeploy} onRemoveProject={handleRemoveProject} />}
      <Route path={`${ROUTES_PREFIX.dashboard}/circles/:circleId/tree`} component={Diagram} />
      <Route path={`${ROUTES_PREFIX.dashboard}/circles/:circleId/tree/deploy`} component={ModalDeploy} />
      <Route path={`${ROUTES_PREFIX.dashboard}/circles/:circleId/tree/add-project`} component={ModalProject} />
    </div>
  )
}

export default Tree