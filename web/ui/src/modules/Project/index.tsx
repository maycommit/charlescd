import React, { useEffect, useState } from 'react'
import { Container } from 'reactstrap'
import { createProject, getProjects } from '../../core/api/project'
import { Button } from 'reactstrap';
import ModalForm from './ModalForm'
import List from './List'
import './style.css'

const Project = () => {
  const [projects, setProjects] = useState<any>([])
  const [modal, setModal] = useState(false);

  const toggle = () => setModal(!modal);

  useEffect(() => {
    (async () => {
      const projects = await getProjects()
      console.log(projects)
      setProjects(projects)
    })()
  }, [])

  const handleSubmit = (data: any) => {
    (async () => {
      createProject(data)
      const projects = await getProjects()
      setProjects(projects)
    })()
  }

  const handleEdit = () => {

  }

  const handleDelete = () => {

  }

  return (
    <Container>
      <div className="header">
        <h3>Projects</h3>
        <Button color="primary" onClick={toggle}>Create</Button>
      </div>
      <List projects={projects} onEdit={handleEdit} onDelete={handleDelete} />
      <ModalForm modal={modal} toggle={toggle} />
    </Container>
  )
}

export default Project