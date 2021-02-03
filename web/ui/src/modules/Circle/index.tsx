import React, { useEffect, useState } from 'react'
import { Link, Route, useHistory, useLocation, useParams, useRouteMatch } from 'react-router-dom'
import { Alert, Badge, Breadcrumb, BreadcrumbItem, Button, CardBody, CardHeader, CardTitle, Col, Container, ListGroup, ListGroupItem, Modal, ModalBody, ModalHeader, Row } from 'reactstrap'
import { createCircle, deleteCircle, deploy, getCircle, getCircles as getCirclesApi, undeploy } from '../../core/api/circle'
import { ROUTES_PREFIX } from '../../core/constants/routes'
import CircleForm from './CircleForm'
import Card from './Card'


const ModalCreate = () => {
  const history = useHistory()
  const match = useRouteMatch()
  const { workspaceId, clusterId } = useParams<any>()

  const toggle = () => history.goBack()

  const handleSubmit = (data: any) => {
    createCircle(workspaceId, clusterId, data)
      .then(() => {
        console.log('Circle CREATE SUCCESS')
      })
    toggle()
  }

  return (
    <Modal isOpen={true} toggle={toggle} className="" size="lg">
      <ModalHeader toggle={toggle}>Create circle</ModalHeader>
      <ModalBody>
        <CircleForm onSubmit={handleSubmit} />
      </ModalBody>
    </Modal>
  )
}



const List = () => {
  const match = useRouteMatch()
  const history = useHistory()
  const location = useLocation()
  const { workspaceId, clusterId } = useParams<any>()
  const [circles, setCircles] = useState([])
  const [error, setError] = useState(null)

  useEffect(() => {
    getCircles()
  }, [location])

  useEffect(() => {
    const interval = setInterval(() => {
      getCircles()
    }, 3000)

    return () => clearInterval(interval)
  }, [])

  const getCircles = () => getCirclesApi(workspaceId, clusterId)
    .then(circles => {
      setCircles(circles)
      setError(null)
    })
    .catch(e => setError(e))

  const handleDeleteCircle = (circleId: string) => deleteCircle(workspaceId, clusterId, circleId)
    .then(() => {
      getCircles()
      setError(null)
    })
    .catch(e => setError(e))

  const getColorByStatus = (status: string) => {
    switch (status) {
      case 'Healthy':
        return 'success'
      case 'Progressing':
        return 'info'
      case 'Hidden':
        return 'light'
      default:
        return 'danger'
    }
  }

  return (
    <div className="cluster__list" style={{ width: '90%', margin: '0 auto' }}>
      <Row className="cluster__list__header" style={{ marginTop: '50px', marginBottom: '40px' }}>
        <Col xs="10">
          <h4 className="text-white">Circles</h4>
        </Col>
        <Col xs="2">
          <Button color="primary" block onClick={() => history.push(`${match.url}/create`)}>Create circle</Button>
        </Col>
      </Row>
      <Row className="mt-3">
        {circles.map((circle: any) => (
          <Col xs="3">
            <Card circle={circle} onCircleClick={(circleName: string) => history.push(`${match.url}/${circleName}/tree`)} />
          </Col>
        ))}
      </Row>

    </div>
  )
}

const Circle = () => {
  const { workspaceId, clusterId } = useParams<any>()
  const match = useRouteMatch()

  return (
    <div className="cluster">
      <Route path={`${ROUTES_PREFIX.dashboard}/circles`} component={List} />
      <Route path={`${ROUTES_PREFIX.dashboard}/circles/create`} component={ModalCreate} />
    </div>
  )

}

export default Circle