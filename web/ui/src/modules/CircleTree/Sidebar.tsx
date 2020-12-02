import React, { useEffect, useState } from 'react'
import AceEditor from "react-ace";

import "ace-builds/src-noconflict/mode-java";
import "ace-builds/src-noconflict/theme-github";
import './style.css'
import { Button } from 'reactstrap';

const Sidebar = ({ circle }: any) => {
  const [name, setName] = useState('')
  const [segments, setSegment] = useState('')
  const [environments, setEnvironments] = useState('')

  useEffect(() => {
    setName(circle?.name)
    setSegment(JSON.stringify(circle?.segments, null, ' '))
    setEnvironments(JSON.stringify(circle?.environments, null, ' '))
  }, [circle])

  return (
    <div className="circle-tree-sidebar">
      <div className="circle-tree-sidebar-header">
        Circle: {name}
      </div>
      <div className="circle-tree-sidebar-section">
        Segments
      </div>
      <AceEditor
        mode="java"
        theme="github"
        onChange={(e: any) => setSegment(e.target.value)}
        value={segments}
        name="UNIQUE_ID_OF_DIV"
        width="100%"
        height="200px"
        editorProps={{ $blockScrolling: true }}
      />
      <div className="circle-tree-sidebar-section">
        Environments
      </div>
      <AceEditor
        mode="java"
        theme="github"
        onChange={(e: any) => setEnvironments(e.target.value)}
        value={environments}
        name="UNIQUE_ID_OF_DIV"
        width="100%"
        height="200px"
        editorProps={{ $blockScrolling: true }}
      />

      <Button color="primary" size="lg" className="circle-tree-sidebar-save" block>Save</Button>
    </div>
  )
}

export default Sidebar