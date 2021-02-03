import React from 'react'
import { NavLink, useLocation, useRouteMatch } from 'react-router-dom'
import { ListGroup } from 'reactstrap'
import './style.scss'

const Sidebar = () => {
  const match = useRouteMatch()
  const location = useLocation()

  return (
    <div className="dashboard__sidebar">
      <ListGroup className="dashboard__sidebar__list">
        <NavLink to={`${match.url}/circles`} className="dashboard__sidebar__list__item" activeClassName="dashboard__sidebar__list__item__active">
          <i className="fas fa-circle-notch"></i>
        </NavLink>
        <NavLink to={`${match.url}/projects`} className="dashboard__sidebar__list__item" activeClassName="dashboard__sidebar__list__item__active">
          <i className="fas fa-folder"></i>
        </NavLink>
      </ListGroup>
    </div>
  )

}

export default Sidebar