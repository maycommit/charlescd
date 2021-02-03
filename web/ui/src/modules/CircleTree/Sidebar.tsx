import React from 'react'
import { useHistory, useParams, useRouteMatch } from 'react-router-dom'
import { Button, Card, CardFooter, CardHeader, ListGroup, ListGroupItem } from 'reactstrap'

const Sidebar = ({ circle, onUndeploy, onRemoveProject }: any) => {
  const history = useHistory()
  const match = useRouteMatch()

  return (
    <div className="tree__sidebar">
      {circle?.release ? (
        <Card style={{ background: '#10AA80' }}>
          <CardHeader>
            <div className="d-flex justify-content-between">
              <div className="text-white">{circle?.release?.name}</div>
              {!circle?.managed && (<i
                className="fas fa-trash"
                onClick={() => window.confirm("Are you sure?") && onUndeploy()}
              />)}
            </div>
          </CardHeader>
          <ListGroup flush>
            {circle?.status?.projects?.map((project: any) => (
              <ListGroupItem className="d-flex justify-content-between text-white" style={{ background: '#10AA80' }}>
                {project?.name}
                {!circle?.managed && (<i
                  className="fas fa-minus-square"
                  onClick={() => window.confirm("Are you sure?") && onRemoveProject(project)}
                />)}
              </ListGroupItem>
            ))}
          </ListGroup>
          <CardFooter>
            {!circle?.managed && <Button block onClick={() => history.push(`${match.url}/add-project`)}>Add project</Button>}
          </CardFooter>
        </Card>
      ) : !circle?.managed && (
        <Button block onClick={() => history.push(`${match.url}/deploy`)}>Add release</Button>
      )}

    </div>
  )
}

export default Sidebar
