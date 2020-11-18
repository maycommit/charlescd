import React, { useState } from 'react'
import { Col, Row, Button, Form, FormGroup, Label, Input, FormText, Modal, ModalHeader, ModalBody, ModalFooter, Alert } from 'reactstrap';

import "ace-builds/src-noconflict/mode-java";
import "ace-builds/src-noconflict/theme-github";

const initialSegment = {
  key: '',
  condition: '',
  value: ''
}

const initialEnvironment = {
  key: '',
  value: ''
}


const DeployForm = ({ onSubmit, toggle, modal, projects, circleName }: any) => {
  const [releaseName, setReleaseName] = useState('')
  const [projectName, setProjectName] = useState('')
  const [tag, setTag] = useState('')

  const handleClick = () => {
    onSubmit(circleName, {
      releaseName,
      projectName,
      tag,
    })
  }

  return (

    <Modal isOpen={modal} toggle={toggle}>
      <ModalHeader toggle={toggle}>Deploy in {circleName}</ModalHeader>
      <ModalBody>
        <Form>
          <FormGroup>
            <Label for="exampleEmail">Release Name</Label>
            <Input type="text" name="name" placeholder="Name..." value={releaseName} onChange={e => setReleaseName(e.target.value)} />
          </FormGroup>

          <FormGroup>
            <Label for="exampleSelect">Projects</Label>
            <Input type="select" name="select" id="exampleSelect" value={projectName} onChange={e => setProjectName(e.target.value)}>
              <option value="" selected disabled hidden>Choose here</option>
              {projects?.map((project: any) => (
                <option value={project?.name}>{project?.name}</option>
              ))}
            </Input>
          </FormGroup>

          <FormGroup>
            <Label for="exampleEmail">Image Tag</Label>
            <Input type="text" name="name" placeholder="Name..." value={tag} onChange={e => setTag(e.target.value)} />
          </FormGroup>

        </Form>
      </ModalBody>
      <ModalFooter>
        <Button color="primary" onClick={handleClick}>Deploy</Button>{' '}
        <Button color="secondary" onClick={toggle}>Cancel</Button>
      </ModalFooter>
    </Modal>
  )
}

export default DeployForm