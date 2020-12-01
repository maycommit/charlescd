import React, { useEffect, useState } from 'react'
import Tree from 'react-d3-tree';
import { useParams } from 'react-router-dom';
import { getCircle } from '../../core/api/circle';
import Sidebar from './Sidebar'
import './style.css'

const CircleTree = () => {
  const { name } = useParams<any>()
  const [circle, setCircle] = useState({})

  useEffect(() => {
    (async () => {
      const circleRes = await getCircle(name)
      setCircle(circleRes)
    })()
  }, [])

  return (
    <>
      <Sidebar
        circle={circle}
      />
      <div className="circle-tree-container">
        <Tree data={{}} />
      </div>
    </>
  )
}

export default CircleTree