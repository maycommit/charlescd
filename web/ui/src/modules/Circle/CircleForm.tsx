import React, { useEffect, useState } from 'react'
import { useFieldArray, useForm } from 'react-hook-form';
import { useParams } from 'react-router-dom';
import { Button, Col, FormGroup, Label, Row } from 'reactstrap'
import { getProject, getProjects } from '../../core/api/project';

const Form = ({ onSubmit }: any) => {
  const { workspaceId, clusterId } = useParams<any>()
  const [projects, setProjects] = useState([])
  const [addEnvironments, setAddEnvironments] = useState(false)
  const [addRelease, setAddRelease] = useState(false)
  const { register, handleSubmit, watch, errors, control } = useForm();
  const { fields: segmentFields, append: appendSegment, remove: removeSegment } = useFieldArray({
    control,
    name: "segments"
  });
  const { fields: environmentFields, append: appendEnvironment, prepend, remove: removeEnvironment } = useFieldArray({
    control,
    name: "environments"
  });
  const onSubmitFunc = (data: any) => onSubmit({ circle: { ...data } });

  const toggleEnvironments = () => {
    setAddEnvironments(!addEnvironments)
  }

  useEffect(() => {
    requestProjects()
    appendSegment({ key: "", condition: "", value: "" })
  }, [])

  useEffect(() => {
    if (addEnvironments) appendEnvironment({ key: "", value: "" })
    else environmentFields.map((e, index) => removeEnvironment(index))
  }, [addEnvironments])

  const requestProjects = () => getProjects(workspaceId, clusterId).then(res => setProjects(res))

  return (
    <div className="cluster">
      <form onSubmit={handleSubmit(onSubmitFunc)}>
        <FormGroup>
          <Label for="name">Name</Label>
          <input name="name" id="name" ref={register({ required: true })} className="form-control" />
          {errors.name && <span>This field is required</span>}
        </FormGroup>
        <Label>Segments</Label>
        <div style={{ border: '1px solid #ccc', padding: '10px', marginBottom: '30px' }}>
          {segmentFields.map((item: any, index: any) => (
            <Row>
              <Col>
                <FormGroup>
                  <Label for="">Key</Label>
                  <input name={`segments[${index}].key`} id="key" ref={register({ required: true })} className="form-control" />
                  {errors.name && <span>This field is required</span>}
                </FormGroup>
              </Col>
              <Col>
                <FormGroup>
                  <Label for="">Condition</Label>
                  <input name={`segments[${index}].condition`} id="condition" ref={register({ required: true })} className="form-control" />
                  {errors.name && <span>This field is required</span>}
                </FormGroup>
              </Col>
              <Col>
                <FormGroup>
                  <Label for="">Value</Label>
                  <input name={`segments[${index}].value`} id="value" ref={register({ required: true })} className="form-control" />
                  {errors.name && <span>This field is required</span>}
                </FormGroup>
              </Col>
              <Col className="d-flex align-items-center">
                <Button type="button" color="primary" onClick={() => appendSegment({ key: "", condition: "", value: "" })}>
                  <i className="fas fa-plus"></i>
                </Button>{' '}
                {segmentFields.length > 1 && (
                  <Button type="button" color="danger" onClick={() => removeSegment(index)}><i className="fas fa-minus"></i></Button>
                )}
              </Col>
            </Row>
          ))}
        </div>
        <Button type="button" onClick={() => toggleEnvironments()} block>
          {addEnvironments ? 'Remove environments' : 'Add environments'}
        </Button>
        {addEnvironments && (
          <div style={{ border: '1px solid #ccc', padding: '10px', marginBottom: '30px' }}>
            {environmentFields.map((item: any, index: any) => (
              <Row>
                <Col>
                  <FormGroup>
                    <Label for="">Key</Label>
                    <input name={`environments[${index}].key`} id="key" ref={register({ required: true })} className="form-control" />
                    {errors.name && <span>This field is required</span>}
                  </FormGroup>
                </Col>
                <Col>
                  <FormGroup>
                    <Label for="">Value</Label>
                    <input name={`environments[${index}].value`} id="value" ref={register({ required: true })} className="form-control" />
                    {errors.name && <span>This field is required</span>}
                  </FormGroup>
                </Col>
                <Col className="d-flex align-items-center">
                  <Button type="button" color="primary" onClick={() => appendEnvironment({ key: "", value: "" })}>
                    <i className="fas fa-plus"></i>
                  </Button>{' '}
                  {environmentFields.length > 1 && (
                    <Button type="button" color="danger" onClick={() => removeEnvironment(index)}><i className="fas fa-minus"></i></Button>
                  )}
                </Col>
              </Row>
            ))}
          </div>
        )}
        <Button type="button" onClick={() => setAddRelease(!addRelease)} block>
          {addRelease ? 'Remove release' : 'Add release'}
        </Button>
        {addRelease && (
          <div style={{ border: '1px solid #ccc', padding: '10px', marginBottom: '30px' }}>
            <FormGroup>
              <Label for="name">Name</Label>
              <input name="release.name" id="name" ref={register({ required: true })} className="form-control" />
              {errors.name && <span>This field is required</span>}
            </FormGroup>
            <FormGroup>
              <Label for="name">Tag</Label>
              <input name="release.tag" id="name" ref={register({ required: true })} className="form-control" />
              {errors.name && <span>This field is required</span>}
            </FormGroup>
            <FormGroup>
              <Label for="projects">Projects</Label>
              {projects.length > 0 && (
                <select name="release.projects" id="projects" ref={register({ required: true })} className="form-control" multiple>
                  {projects.map((project: any, index: any) => (
                    <option value={project?.name}>{project?.name}</option>
                  ))}
                </select>
              )}

              {errors.name && <span>This field is required</span>}
            </FormGroup>

          </div>
        )}
        <FormGroup>
          <Label for="">Namespace</Label>
          <input name={`destination.namespace`} id="key" ref={register({ required: true })} className="form-control" value="default" readOnly />
          {errors.name && <span>This field is required</span>}
        </FormGroup>
        <Button color="primary" className="mt-4">Create</Button>
      </form>
    </div >
  )

}

export default Form