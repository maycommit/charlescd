import React, { useState } from 'react'
import { Col, Row, Button, Form, FormGroup, Label, Input, FormText, Modal, ModalHeader, ModalBody, ModalFooter } from 'reactstrap';
import { createProject } from '../../core/api/project';

const ProjectForm = ({ onSubmit, toggle, modal }: any) => {
  const [name, setName] = useState('')
  const [repository, setRepository] = useState('')
  const [paths, setPaths] = useState('')
  const [token, setToken] = useState('')

  const splitPaths = (pathsStr: string) => {
    return pathsStr.split(';')
  }

  const handleClick = () => {
    onSubmit({
      name,
      repository,
      paths: splitPaths(paths),
      token
    })
  }

  return (

    <Modal isOpen={modal} toggle={toggle}>
      <ModalHeader toggle={toggle}>Create project</ModalHeader>
      <ModalBody>
        <Form>
          <FormGroup>
            <Label for="exampleEmail">Name</Label>
            <Input type="text" name="name" placeholder="Name..." value={name} onChange={e => setName(e.target.value)} />
          </FormGroup>

          <FormGroup>
            <Label for="exampleEmail">Repository</Label>
            <Input type="text" name="repository" placeholder="Repository..." value={repository} onChange={e => setRepository(e.target.value)} />
          </FormGroup>

          <FormGroup>
            <Label for="exampleEmail">Paths</Label>
            <Input type="text" name="paths" placeholder="Paths..." value={paths} onChange={e => setPaths(e.target.value)} />
            <FormText>Separated by ';'</FormText>
          </FormGroup>

          <FormGroup>
            <Label for="exampleEmail">Token</Label>
            <Input type="text" name="token" placeholder="Token..." value={token} onChange={e => setToken(e.target.value)} />
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

export default ProjectForm