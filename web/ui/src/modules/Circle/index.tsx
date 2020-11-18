import React, { useEffect, useState } from 'react'
import { Container, Button } from 'reactstrap'
import ModalForm from './ModalForm'
import List from './List'
import DeployModalForm from './ModalDeploy'
import { createCircle, deploy, getCircles } from '../../core/api/circle'
import { getProjects } from '../../core/api/project'

const Circle = () => {
  const [modal, setModal] = useState(false);
  const [modalDeploy, setModalDeploy] = useState(false);
  const [circleSelected, setCircleSelected] = useState('');
  const [circles, setCircles] = useState([])
  const [projects, setProjects] = useState([])

  const toggle = () => setModal(!modal);

  const toggleModalDeploy = () => setModalDeploy(!modalDeploy);

  const handleOnDeploy = (name: string) => {
    toggleModalDeploy()
    setCircleSelected(name)
  }

  const handleOnDeploySubmit = (name: string, data: any) => {
    (async () => {
      await deploy(name, data)
    })()
  }

  useEffect(() => {
    (async () => {
      try {
        const projects = await getProjects()
        const circles = await getCircles()
        setCircles(circles)
        setProjects(projects)
      } catch (e) {
        console.error(e)
      }
    })()
  }, [])

  const handleSubmit = (data: any) => {
    (async () => {
      try {
        createCircle(data)
        const circles = await getCircles()
        setCircles(circles)
      } catch (e) {
        console.error(e)
      }
    })()
  }

  return (
    <Container>
      <div className="header">
        <h3>Circles</h3>
        <Button color="primary" onClick={toggle}>Create</Button>
      </div>
      <List circles={circles} onDeploy={handleOnDeploy} />
      {/* <List projects={projects} onEdit={handleEdit} onDelete={handleDelete} /> */}
      <ModalForm modal={modal} toggle={toggle} onSubmit={handleSubmit} />
      <DeployModalForm modal={modalDeploy} toggle={toggleModalDeploy} projects={projects} onSubmit={handleOnDeploySubmit} circleName={circleSelected} />
    </Container>
  )
}

export default Circle