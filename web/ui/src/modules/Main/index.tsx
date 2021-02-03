import React from 'react'
import { Switch, Route, Redirect } from 'react-router-dom'
import Sidebar from './Sidebar'
import Workspace from '../Workspace'

import './style.scss'

const Main = () => {

  return (
    <div className="main">
      <Sidebar />

      <div className="main__content">
        <Switch>
          <Redirect exact path="/" to="/workspaces" />
          <Route path="/workspaces" component={Workspace} />
        </Switch>
      </div>
    </div>
  )
}

export default Main