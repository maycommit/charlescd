import React, { useEffect, useState } from 'react'
import { Col, Row, Button, Form, FormGroup, Label, Input, FormText, Modal, ModalHeader, ModalBody, ModalFooter, Alert } from 'reactstrap';
import { createCircle } from '../../core/api/circle';

const CircleForm = ({ onSubmit, toggle, circle, modal }: any) => {
  const [name, setName] = useState('')

  useEffect(() => {
    setName(circle?.name)
  }, [circle])

  const handleClick = () => {
    try {
      onSubmit({
        name,
      })
    } catch (e) {

    }
  }

  return (
    <Modal isOpen={modal} toggle={toggle}>
      <ModalHeader toggle={toggle}>Create circle</ModalHeader>
      <ModalBody>
        <Form>
          <FormGroup>
            <Label for="exampleEmail">Name</Label>
            <Input type="text" name="name" placeholder="Name..." value={name} onChange={e => setName(e.target.value)} />
          </FormGroup>
        </Form>
      </ModalBody>
      <ModalFooter>
        <Button color="primary" onClick={handleClick}>Save</Button>{' '}
        <Button color="secondary" onClick={toggle}>Cancel</Button>
      </ModalFooter>
    </Modal>
  )
}

export default CircleForm