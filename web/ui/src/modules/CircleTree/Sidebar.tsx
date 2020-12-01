import React, { useEffect, useState } from 'react'
import './style.css'

const Sidebar = ({ circle }: any) => {
  const [name, setName] = useState('')

  useEffect(() => {
    setName(circle?.name)
  }, [circle])

  return (
    <div className="circle-tree-sidebar">
      <h3>Circle: {name}</h3>
    </div>
  )
}

export default Sidebar