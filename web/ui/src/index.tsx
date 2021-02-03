import React from 'react';
import ReactDOM from 'react-dom';
import {
  BrowserRouter as Router,
  Switch,
  Route,
  Link
} from "react-router-dom";
import Main from './modules/Main';
import reportWebVitals from './reportWebVitals';
import 'bootstrap/dist/css/bootstrap.min.css';
import './index.scss';
import Dashboard from './modules/Dashboard';
import { ROUTES_PREFIX } from './core/constants/routes';
import './index.scss'

ReactDOM.render(
  <React.StrictMode>
    <Router>
      <Switch>
        <Route path={ROUTES_PREFIX.dashboard} component={Dashboard} />
        <Route path={ROUTES_PREFIX.main} component={Main} />
      </Switch>
    </Router>
  </React.StrictMode>,
  document.getElementById('root')
);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
