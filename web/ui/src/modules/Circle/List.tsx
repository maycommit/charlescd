import React, { useState } from 'react'
import { Link } from 'react-router-dom';
import { Col, Row, Popover, PopoverHeader, PopoverBody, ListGroup, ListGroupItem } from 'reactstrap'
import {
  Card, Button, CardHeader, CardFooter, CardBody,
  CardTitle, CardText
} from 'reactstrap';
import './style.css'

const PopoverItem = ({ id, title, children }: any) => {
  const [popoverOpen, setPopoverOpen] = useState(false);

  const toggle = () => setPopoverOpen(!popoverOpen);

  return (
    <>
      <Button id={`popover-${id}`} block>{title}</Button>
      <Popover placement="right" isOpen={popoverOpen} target={`popover-${id}`} toggle={toggle}>
        <PopoverHeader>{title}</PopoverHeader>
        <PopoverBody>{children}</PopoverBody>
      </Popover>
    </>
  )
}

const List = ({ circles, onDeploy, onEdit, onDelete }: any) => {
  return (
    <Row>
      {circles.map((circle: any) => (
        <Col xs={3}>
          <Card>
            <CardHeader>
              <Link to={`/circles/${circle?.name}`}>{circle?.name}</Link>
              <div className="card-header-icons">
                <a onClick={() => onEdit(circle)}><i className="fas fa-pen"></i></a>
                <a onClick={() => onDelete(circle)}><i className="fas fa-trash"></i></a>
              </div>
            </CardHeader>
            <CardBody>

            </CardBody>
            <CardFooter>

            </CardFooter>
          </Card>
        </Col>
      ))}
    </Row>
  )
}

export default List