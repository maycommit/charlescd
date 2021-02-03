import React, { useEffect } from 'react'
import { useForm } from 'react-hook-form';
import { ModalFooter, Button, FormGroup, Input, Label } from 'reactstrap';
import './style.scss'

const Form = ({ onSubmit }: any) => {
  const { register, handleSubmit, watch, errors } = useForm();
  const onSubmitFunc = (data: any) => onSubmit({ ...data });

  return (
    <div className="workspace__form">
      <form onSubmit={handleSubmit(onSubmitFunc)}>
        <FormGroup>
          <input name="name" id="name" ref={register({ required: true })} className="form-control" placeholder="Name" />
          {errors.name && <span>This field is required</span>}
        </FormGroup>
        <FormGroup>
          <textarea name="description" id="description" ref={register} className="form-control" placeholder="Description"></textarea>
          {errors.name && <span>This field is required</span>}
        </FormGroup>

        <Button type="submit" color="primary">Create</Button>{' '}
      </form>
    </div >
  )
}

export default Form