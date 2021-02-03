import React, { useEffect, useState } from 'react'
import { useParams } from 'react-router-dom'
import { Alert, Card, CardBody, CardHeader, Col, Jumbotron, Row } from 'reactstrap'
import dayjs from 'dayjs'
import { getManifest } from '../../core/api/circle'
import AceEditor from "react-ace";
import { getColorByHealth } from './helper'
import './style.scss'

import "ace-builds/src-noconflict/mode-json";
import "ace-builds/src-noconflict/theme-monokai";

const Sidebar = ({ node, onClose }: any) => {
  const [manifest, setManifest] = useState<any>(null)
  const { workspaceId, clusterId, circleId } = useParams<any>()
  const getHealthColor: any = {
    "Healthy": "success",
    "Degraded": 'danger'
  }

  useEffect(() => {
    const group = node?.meta?.group || ''
    const version = node?.meta?.version || ''

    getManifest(workspaceId, clusterId, node?.meta?.name, version, group, node?.meta?.kind)
      .then((data: any) => setManifest(data))
      .catch(() => setManifest(null))
  }, [node])

  return (
    <div className="tree__node-sidebar">
      <Row>
        <Col xs="11">
          <h5 className="mb-3">{node?.name}</h5>
        </Col>
        <Col xs="1">
          <i className="fas fa-times" style={{ cursor: 'pointer' }} onClick={onClose}></i>
        </Col>
      </Row>
      <hr className="my-2" />
      <div className="d-flex justify-content-between mb-4">
        <span><strong>Kind: </strong>{node?.meta?.kind}</span>
        <span><strong>Creation: </strong>{dayjs(node.creationTime).format('DD-MM-YYYY')} - {dayjs(node.creationTime).format('HH:mm')}</span>
      </div>
      {node?.meta?.health && node?.meta?.health?.status && (
        <div className="tree__node-sidebar__node">
          <Card className={`tree__node-sidebar__node__${node?.meta?.health?.status}`}>
            <div><strong>{node?.meta?.health?.status}</strong></div>
            {node?.meta?.health.message && <p>{node?.meta?.health.message}</p>}
          </Card>
        </div>
      )}

      {manifest && (
        <AceEditor
          style={{ borderRadius: '10px', border: '2px solid #212121' }}
          mode="json"
          readOnly={true}
          theme="monokai"
          value={JSON.stringify(manifest, null, ' ')}
          name="manifest"
          height="300px"
          width="100%"
        />
      )}
    </div>
  )
}

export default Sidebar