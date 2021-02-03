import { cluster, pointer } from 'd3'
import React from 'react'
import { Button, Card as ReactstrapCard, CardBody, CardSubtitle, CardTitle, ListGroup, ListGroupItem } from 'reactstrap'
import './style.scss'

const Card = ({ workspace, onClusterClick, onClusterCreate }: any) => {

  return (
    <div className="workspace__card">
      <div>{workspace?.name}</div>
      <small>{workspace?.description}</small>

      <div className="workspace__card__clusters">
        {workspace?.clusters?.map((cluster: any) => (
          <div className="workspace__card__clusters__item" onClick={() => onClusterClick(workspace?.id, cluster?.id)}>
            <div>{cluster?.name}</div>
            <small>{cluster?.address}</small>
          </div>
        ))}
        {workspace?.clusters?.length <= 0 && (
          <Button block color="primary" onClick={() => onClusterCreate(workspace?.id)}>Create cluster</Button>
        )}
      </div>
    </div>
  )
}

export default Card