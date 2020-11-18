import React, { useState } from 'react'
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

const List = ({ circles, onDeploy }: any) => {
  return (
    <Row>
      {circles.map((circle: any) => (
        <Col xs={3}>
          <Card>
            <CardHeader>{circle?.name}</CardHeader>
            <CardBody>
              <PopoverItem id="segments" title="Segments">
                <ListGroup>
                  {circle?.segments?.map((segment: any) => (
                    <ListGroupItem>
                      {`${segment?.key}`}{' '}
                      {`${segment?.condition}`}{' '}
                      {`${segment?.value}`}
                    </ListGroupItem>
                  ))}
                </ListGroup>
              </PopoverItem>
              <PopoverItem id="environments" title="Environments">
                <ListGroup>
                  {circle?.environments?.map((environment: any) => (
                    <ListGroupItem>
                      <div><strong>Key: </strong>{`${environment?.key}`}{' '}</div>
                      <div><strong>Value: </strong>{`${environment?.value}`}{' '}</div>
                    </ListGroupItem>
                  ))}
                </ListGroup>
              </PopoverItem>
              <div>
                {circle?.resources?.map((resource: any) => (
                  <div className="circle-resource">
                    <div>{resource?.name}</div>
                  </div>
                ))}
              </div>
            </CardBody>
            <CardFooter>
              {circle?.resources?.length <= 0 && <Button color="primary" onClick={() => onDeploy(circle?.name)} block>Deploy</Button>}
            </CardFooter>
          </Card>
        </Col>
      ))}
    </Row>
  )
}

export default List