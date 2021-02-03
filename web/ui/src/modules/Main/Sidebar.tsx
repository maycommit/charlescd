import React from 'react'
import { NavLink, useLocation, useRouteMatch } from 'react-router-dom'
import { ListGroup } from 'reactstrap'
import './style.scss'

const Sidebar = () => {
  const match = useRouteMatch()
  const location = useLocation()

  return (
    <div className="main__sidebar">
      <ListGroup className="main__sidebar__list">
        <NavLink to={`/workspaces`} className="main__sidebar__list__item" activeClassName="main__sidebar__list__item__active">
          <i className="fas fa-border-all"></i>
        </NavLink>
      </ListGroup>
    </div>
  )

}

export default Sidebar