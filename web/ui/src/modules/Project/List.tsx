import React, { useEffect, useState } from 'react'
import { Table, Button } from 'reactstrap';
import './style.css'

const List = ({ projects, onEdit, onDelete }: any) => {
  return (
    <Table>
      <thead>
        <tr>
          <th>Name</th>
          <th>Repository</th>
          <th>Paths</th>
          <th>Actions</th>
        </tr>
      </thead>
      <tbody>
        {projects?.map((project: any) => (
          <tr key={project?.name}>
            <th scope="row">{project?.name}</th>
            <td>{project?.repository}</td>
            <td>{project?.paths}</td>
            <td>
              <Button color="primary" onClick={() => onEdit(project)}>Edit</Button>{' '}
              <Button color="danger" onClick={() => onDelete(project?.name)}>Delete</Button>
            </td>
          </tr>
        ))}
      </tbody>
    </Table>
  )
}

export default List