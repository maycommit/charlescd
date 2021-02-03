import React, { useState } from 'react'
import { Modal, ModalHeader, ModalBody } from 'reactstrap';
import { createWorkspace } from '../../core/api/workspace'
import WorkspaceForm from './WorkspaceForm';
import ClusterForm from './ClusterForm'
import List from './List'
import './style.scss'
import { createCluster } from '../../core/api/cluster';
import { Link, Route, useHistory, useParams } from 'react-router-dom';

const ModalCreateWorkspace = () => {
  const history = useHistory()

  const toggle = () => history.push('/workspaces')

  const handleSubmit = (data: any) =>
    createWorkspace(data)
      .then(toggle)

  return (
    <Modal isOpen={true} toggle={toggle} className="">
      <ModalHeader toggle={toggle}>Create Workspace</ModalHeader>
      <ModalBody>
        <WorkspaceForm onSubmit={handleSubmit} />
      </ModalBody>
    </Modal>
  )
}


const ModalCreateCluster = () => {
  const history = useHistory()
  const { workspaceId } = useParams<any>()

  const toggle = () => history.push('/workspaces')

  const handleSubmit = (data: any) =>
    createCluster(workspaceId, data)
      .then(toggle)

  return (
    <Modal isOpen={true} toggle={toggle} className="">
      <ModalHeader toggle={toggle}>Create cluster</ModalHeader>
      <ModalBody>
        <ClusterForm onSubmit={handleSubmit} />
      </ModalBody>
    </Modal>
  )
}

const Workspace = () => {
  return (
    <div className="workspace">
      <div className="workspace__list">
        <Route path="/workspaces" component={List} />
        <Route path="/workspaces/create" component={ModalCreateWorkspace} />
        <Route path="/workspaces/:workspaceId/create-cluster" component={ModalCreateCluster} />
      </div>
    </div >
  )
}

export default Workspace