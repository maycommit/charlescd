import React from 'react'
import {
  BrowserRouter as Router,
  Switch,
  Route,
  Link
} from "react-router-dom";
import Circle from '../Circle';
import Project from '../Project';
import Navbar from './Navbar'

const App = () => {
  return (
    <Router>
      <div>
        <Navbar />

        <div>
          <Switch>
            <Route path="/circles" component={Circle} />
            <Route path="/projects" component={Project} />
          </Switch>
        </div>
      </div>
    </Router>
  )
}

export default App