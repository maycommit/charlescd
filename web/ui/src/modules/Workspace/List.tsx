import React, { useEffect, useState } from 'react'
import { Button, ListGroup, ListGroupItem, ListGroupItemHeading, ListGroupItemText, Row, Col, Container } from 'reactstrap';
import { deleteWorkspace, getWorkspaces as getWorkspacesApi } from '../../core/api/workspace'
import './style.scss'
import { deleteCluster } from '../../core/api/cluster';
import { Link, useHistory, useLocation } from 'react-router-dom';
import Card from './Card'

const List = () => {
  const history = useHistory()
  const location = useLocation()
  const [workspaces, setWorkspaces] = useState<any>([])

  useEffect(() => {
    getWorkspaces()
  }, [location])

  const getWorkspaces = () => getWorkspacesApi()
    .then(res => setWorkspaces(res))

  const handleDeleteWorkspace = (workspaceId: string) => window.confirm('Are you sure?') && deleteWorkspace(workspaceId)
    .then(() => getWorkspaces())

  const handleDeleteCluster = (workspaceId: string, clusterId: string) => window.confirm('Are you sure?') && deleteCluster(workspaceId, clusterId).then(() => getWorkspaces())

  return (
    <div style={{ width: '90%', margin: '0 auto' }}>
      <Row className="cluster__list__header" style={{ marginTop: '50px', marginBottom: '40px' }}>
        <Col xs="10">
          <h4 className="text-white">Workspaces</h4>
        </Col>
        <Col xs="2">
          <Button color="primary" block onClick={() => history.push(`/workspaces/create`)}>Create workspace</Button>
        </Col>
      </Row>

      <Row className="mt-4">
        {workspaces?.map((workspace: any) => (
          <Col xs="3">
            <Card
              workspace={workspace}
              onClusterClick={(workspaceId: string, clusterId: string) => history.push(`/workspaces/${workspaceId}/clusters/${clusterId}/circles`)}
              onClusterCreate={(workspaceId: string) => history.push(`/workspaces/${workspaceId}/create-cluster`)}
            />
          </Col>
        ))}
      </Row>
    </div>
  )
}

export default List