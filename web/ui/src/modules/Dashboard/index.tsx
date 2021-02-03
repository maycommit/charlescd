import React, { useEffect, useState } from 'react'
import { Link, Route, Switch, useLocation, useParams, useRouteMatch } from 'react-router-dom'
import { Alert, Breadcrumb, BreadcrumbItem } from 'reactstrap'
import { clusterHealth } from '../../core/api/cluster'
import { ROUTES_PREFIX } from '../../core/constants/routes'
import Circle from '../Circle'
import CircleTree from '../CircleTree'
import Cluster from '../Cluster'
import Project from '../Project'
import Sidebar from './Sidebar'
import './style.scss'

const Dashboard = () => {
  const { workspaceId, clusterId } = useParams<any>()
  const [error, setError] = useState<any>(null)

  return (
    <div className="dashboard">
      <Sidebar />

      <div className="dashboard__content">
        {error && (<Alert color="danger">
          {error?.message}
        </Alert>)}
        <Switch>
          <Route path={`${ROUTES_PREFIX.dashboard}`} component={Cluster} exact />
          <Route path={`${ROUTES_PREFIX.dashboard}/circles/:circleId/tree`} component={CircleTree} />
          <Route path={`${ROUTES_PREFIX.dashboard}/circles`} component={Circle} />
          <Route path={`${ROUTES_PREFIX.dashboard}/projects`} component={Project} />
        </Switch>
      </div>
    </div>
  )

}

export default Dashboard