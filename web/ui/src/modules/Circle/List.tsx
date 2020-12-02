import React from 'react'
import { Link } from 'react-router-dom';
import { Col, Row } from 'reactstrap'
import {
  Card, CardHeader, CardFooter, CardBody,
} from 'reactstrap';
import './style.css'

// const PopoverItem = ({ id, title, children }: any) => {
//   const [popoverOpen, setPopoverOpen] = useState(false);

//   const toggle = () => setPopoverOpen(!popoverOpen);

//   return (
//     <>
//       <Button id={`popover-${id}`} block>{title}</Button>
//       <Popover placement="right" isOpen={popoverOpen} target={`popover-${id}`} toggle={toggle}>
//         <PopoverHeader>{title}</PopoverHeader>
//         <PopoverBody>{children}</PopoverBody>
//       </Popover>
//     </>
//   )
// }

const List = ({ circles, onDeploy, onEdit, onDelete }: any) => {
  return (
    <Row>
      {circles.map((circle: any) => (
        <Col id={circle?.id} xs={3}>
          <Card>
            <CardHeader>
              <Link id={circle?.name} to={`/circles/${circle?.name}`}>{circle?.name}</Link>
              <div className="card-header-icons">
                <span onClick={() => onEdit(circle)}><i className="fas fa-pen"></i></span>
                <span onClick={() => onDelete(circle)}><i className="fas fa-trash"></i></span>
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