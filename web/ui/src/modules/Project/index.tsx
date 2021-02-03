import React, { useEffect, useState } from 'react'
import { Link, Route, useHistory, useLocation, useParams, useRouteMatch } from 'react-router-dom'
import { Alert, Breadcrumb, BreadcrumbItem, Button, Card, CardBody, CardHeader, CardSubtitle, CardTitle, Col, Container, ListGroup, ListGroupItem, Modal, ModalBody, ModalHeader, Row, Table, UncontrolledCollapse } from 'reactstrap'
import { createProject, deleteProject, getProjects as getProjectsApi } from '../../core/api/project'
import { ROUTES_PREFIX } from '../../core/constants/routes'
import ProjectForm from './ProjectForm'


const ModalCreate = () => {
  const history = useHistory()

  const match = useRouteMatch()
  const { workspaceId, clusterId } = useParams<any>()

  const toggle = () => history.goBack()

  const handleSubmit = (data: any) => {
    createProject(workspaceId, clusterId, data)
      .then(() => {
        console.log('Project CREATE SUCCESS')
      })
    toggle()
  }

  return (
    <Modal isOpen={true} toggle={toggle} className="" size="lg">
      <ModalHeader toggle={toggle}>Create project</ModalHeader>
      <ModalBody>
        <ProjectForm onSubmit={handleSubmit} />
      </ModalBody>
    </Modal>
  )
}

const List = () => {
  const match = useRouteMatch()
  const history = useHistory()
  const location = useLocation()
  const { workspaceId, clusterId } = useParams<any>()
  const [projects, setProjects] = useState([])

  useEffect(() => {
    getProjects()
  }, [location])

  useEffect(() => {
    const interval = setInterval(() => {
      getProjects()
    }, 3000)


    return () => clearInterval(interval)
  }, [])

  const getProjects = () => getProjectsApi(workspaceId, clusterId)
    .then(projects => setProjects(projects))

  const handleDeleteProject = (circleid: string) => deleteProject(workspaceId, clusterId, circleid)
    .then(() => getProjects())

  return (
    <div style={{ width: '90%', margin: '0 auto' }}>
      <Row className="cluster__list__header" style={{ marginTop: '50px', marginBottom: '40px' }}>
        <Col xs="10">
          <h4 className="text-white">Projects</h4>
        </Col>
        <Col xs="2">
          <Button color="primary" block onClick={() => history.push(`${match.url}/create`)}>Create project</Button>
        </Col>
      </Row>
      <hr />
      <Row>
        {projects.map((project: any, i: any) => (
          <Col xs="3">
            <Card className="mb-3" style={{ padding: '10px', background: '#3A3A3C', color: '#fff' }}>
              <Row>
                <Col xs="10">
                  <div>{project?.name}</div>
                  <div className="mb-2 text-muted">{project?.repoUrl}/{project?.path}</div>
                </Col>
                <Col xs="2" className="d-flex align-items-center">
                  <Button block color="link" id={`toggler-${i}`} style={{ outline: 'none' }}>
                    <i className="fas fa-project-diagram"></i>
                  </Button>
                </Col>
              </Row>
              <UncontrolledCollapse toggler={`#toggler-${i}`}>
                {project?.routes?.length <= 0 && (
                  <div className="text-white p-3">No routes for project</div>
                )}
                {project?.routes?.map((route: any) => (
                  <Card
                    style={{ background: ' #10AA80', padding: '5px', cursor: 'pointer' }}
                    onClick={() => history.push(`/workspaces/${workspaceId}/clusters/${clusterId}/circles/${route?.circleName}/tree`)}
                  >
                    <div><strong>Circle:</strong> {route?.circleName}</div>
                    <div><strong>Release:</strong> {route?.releaseName}</div>
                  </Card>
                ))}
              </UncontrolledCollapse>
            </Card>
          </Col>
        ))}
      </Row>

    </div>
  )
}

const Project = () => {
  const { workspaceId, clusterId } = useParams<any>()
  return (
    <div className="project">
      <Route path={`${ROUTES_PREFIX.dashboard}/projects`} component={List} />
      <Route path={`${ROUTES_PREFIX.dashboard}/projects/create`} component={ModalCreate} />
    </div>
  )

}

export default Project