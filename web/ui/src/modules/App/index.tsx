import React from 'react'
import {
  BrowserRouter as Router,
  Switch,
  Route,
  Link,
  Redirect
} from "react-router-dom";
import Circle from '../Circle';
import CircleTree from '../CircleTree';
import Navbar from './Navbar'

const App = () => {
  return (
    <Router>
      <div>
        <Navbar />

        <div>
          <Switch>
            <Route path="/" component={Circle} exact />
            <Route path="/circles" component={Circle} exact />
            <Route path="/circles/:name" component={CircleTree} />
          </Switch>
        </div>
      </div>
    </Router>
  )
}

export default App