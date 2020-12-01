import React, { useEffect, useState } from 'react'
import { Container, Button } from 'reactstrap'
import ModalForm from './ModalForm'
import List from './List'
import { createCircle, deploy, getCircles } from '../../core/api/circle'
import { getProjects } from '../../core/api/project'

const Circle = () => {
  const [modal, setModal] = useState(false);
  const [circles, setCircles] = useState([])
  const [circle, setCircle] = useState(null)

  const toggle = () => {
    setCircle(null)
    setModal(!modal);
  }

  useEffect(() => {
    (async () => {
      try {
        const circles = await getCircles()
        setCircles(circles)
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

  const handleEdit = (circle: any) => {
    setCircle(circle)
    setModal(!false)
  }

  const handleDelete = (circle: any) => {
    if (window.confirm("Do you really want to delete?")) {
      console.log("Confirm delete")
      return
    }

    console.log('No delete')
  }

  return (
    <Container>
      <div className="header">
        <h3>Circles</h3>
        <Button color="primary" onClick={toggle}>Create</Button>
      </div>
      <List
        circles={circles}
        onEdit={handleEdit}
        onDelete={handleDelete}
      />
      {/* <List projects={projects} onEdit={handleEdit} onDelete={handleDelete} /> */}
      <ModalForm
        modal={modal}
        circle={circle}
        toggle={toggle}
        onSubmit={handleSubmit}
      />
    </Container>
  )
}

export default Circle