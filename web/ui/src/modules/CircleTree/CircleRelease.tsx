import React, { useEffect, useState } from 'react'
import { useFieldArray, useForm } from 'react-hook-form';
import { useParams } from 'react-router-dom';
import { Button, Col, FormGroup, Label, Row } from 'reactstrap'
import { getProject, getProjects } from '../../core/api/project';

const Form = ({ onSubmit }: any) => {
  const { workspaceId, clusterId } = useParams<any>()
  const [projects, setProjects] = useState([])
  const { register, handleSubmit, watch, errors, control } = useForm();
  const onSubmitFunc = (data: any) => onSubmit({ release: { ...data } });

  useEffect(() => {
    requestProjects()
  }, [])

  const requestProjects = () => getProjects(workspaceId, clusterId).then(res => setProjects(res))

  return (
    <div className="cluster">
      <form onSubmit={handleSubmit(onSubmitFunc)}>
        <FormGroup>
          <Label for="name">Name</Label>
          <input name="name" id="name" ref={register({ required: true })} className="form-control" />
          {errors.name && <span>This field is required</span>}
        </FormGroup>
        <FormGroup>
          <Label for="name">Tag</Label>
          <input name="tag" id="name" ref={register({ required: true })} className="form-control" />
          {errors.name && <span>This field is required</span>}
        </FormGroup>
        <FormGroup>
          <Label for="projects">Projects</Label>
          {projects.length > 0 && (
            <select name="projects" id="projects" ref={register({ required: true })} className="form-control" multiple>
              {projects.map((project: any, index: any) => (
                <option value={project?.name}>{project?.name}</option>
              ))}
            </select>
          )}

          {errors.name && <span>This field is required</span>}
        </FormGroup>
        <Button color="primary" className="mt-4">Deploy</Button>
      </form>
    </div >
  )

}

export default Form