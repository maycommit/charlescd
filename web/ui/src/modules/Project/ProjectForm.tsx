import React, { useEffect, useState } from 'react'
import { useFieldArray, useForm } from 'react-hook-form';
import { Button, Col, FormGroup, Label, Row } from 'reactstrap'

const Form = ({ onSubmit }: any) => {
  const { register, handleSubmit, watch, errors, control } = useForm();
  const onSubmitFunc = (data: any) => onSubmit({ project: { ...data } });


  return (
    <div className="cluster">
      <form onSubmit={handleSubmit(onSubmitFunc)}>
        <FormGroup>
          <Label for="name">Name</Label>
          <input name="name" id="name" ref={register({ required: true })} className="form-control" />
          {errors.name && <span>This field is required</span>}
        </FormGroup>
        <FormGroup>
          <Label for="name">Repository Url</Label>
          <input name="repoUrl" id="name" ref={register({ required: true })} className="form-control" />
          {errors.name && <span>This field is required</span>}
        </FormGroup>
        <FormGroup>
          <Label for="name">Path</Label>
          <input name="path" id="name" ref={register} className="form-control" />
          {errors.name && <span>This field is required</span>}
        </FormGroup>
        <FormGroup>
          <Label for="name">Token</Label>
          <input name="token" id="name" ref={register} className="form-control" />
          {errors.name && <span>This field is required</span>}
        </FormGroup>
        <FormGroup>
          <Label for="name">Template type</Label>
          <select name="template.type" id="name" ref={register} className="form-control">
            <option value="puremanifest">Pure manifest</option>
          </select>
          {errors.name && <span>This field is required</span>}
        </FormGroup>
        <Button color="primary">Create</Button>
      </form>
    </div >
  )

}

export default Form