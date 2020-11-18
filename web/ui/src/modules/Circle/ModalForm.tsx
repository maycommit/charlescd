import React, { useState } from 'react'
import { Col, Row, Button, Form, FormGroup, Label, Input, FormText, Modal, ModalHeader, ModalBody, ModalFooter, Alert } from 'reactstrap';
import AceEditor from "react-ace";
import { createCircle } from '../../core/api/circle';

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


const CircleForm = ({ onSubmit, toggle, modal }: any) => {
  const [addEnvironment, setAddEnvironment] = useState(false)
  const [name, setName] = useState('')
  const [segments, setSegments] = useState(JSON.stringify([initialSegment], null, 2))
  const [environments, setEnvironments] = useState<any>(JSON.stringify([initialEnvironment], null, 2))
  const [error, setError] = useState('')

  const handleClick = () => {
    try {
      onSubmit({
        name,
        segments: JSON.parse(segments),
        environments: JSON.parse(environments),
        resources: [],
      })
    } catch (e) {
      setError(e)
    }
  }

  return (

    <Modal isOpen={modal} toggle={toggle}>
      <ModalHeader toggle={toggle}>Create circle</ModalHeader>
      <ModalBody>
        <Alert color="danger" isOpen={error !== ''} toggle={() => setError('')}>
          {error}
        </Alert>
        <Form>
          <FormGroup>
            <Label for="exampleEmail">Name</Label>
            <Input type="text" name="name" placeholder="Name..." value={name} onChange={e => setName(e.target.value)} />
          </FormGroup>

          <FormGroup>
            <Label for="exampleEmail">Segments</Label>
            <AceEditor
              mode="json"
              theme="github"
              onChange={values => setSegments(values)}
              value={segments}
              name="segments"
              width='100%'
              height="200px"
              editorProps={{ $blockScrolling: true }}
            />
          </FormGroup>


          <FormGroup>
            <Label for="exampleEmail">Environments</Label>
            <AceEditor
              mode="json"
              theme="github"
              onChange={values => setEnvironments(values)}
              value={environments}
              name="enviroments"
              width='100%'
              height="200px"
              editorProps={{ $blockScrolling: true }}
            />
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