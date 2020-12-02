import React from 'react'

const Resource = ({ node }: any) => {
  const health = node.meta.health
  let healthColor = "green"

  if (health && health?.status !== "Healthy") {
    healthColor = "red"
  } else if (!health) {
    healthColor = "grey"
  }

  return (
    <div
      style={{
        borderRadius: "5px",
        borderWidth: "2px 2px 2px 10px",
        borderColor: healthColor,
        borderStyle: "solid",
        maxWidth: "200px",
        minWidth: "180px",
        padding: "10px 10px",
        backgroundColor: "#fff"
      }}
    >
      <div style={{ fontWeight: "bold" }}>{node.label}</div>
      <div style={{ fontSize: "10px" }}>{node.meta.creationTimestamp}</div>

      {health && health?.status !== "Healthy" && (<div className="circle-tree-node-error">{health?.message}</div>)}
    </div>
  )
}

const Project = ({ node }: any) => {
  return (
    <div
      style={{
        borderRadius: "5px",
        borderWidth: "2px 2px 2px 10px",
        borderColor: "black",
        borderStyle: "solid",
        maxWidth: "200px",
        minWidth: "180px",
        padding: "10px 10px",
        backgroundColor: "#fff"
      }}
    >
      <div style={{ fontWeight: "bold" }}>{node.label}</div>
    </div>
  )
}

const Circle = ({ node }: any) => {
  return (
    <div
      style={{
        borderRadius: "5px",
        borderWidth: "2px",
        borderColor: "#000",
        borderStyle: "solid",
        maxWidth: "200px",
        minWidth: "180px",
        padding: "10px 10px",
        backgroundColor: "#fff"
      }}
    >
      <div style={{ fontWeight: "bold" }}>{node.label}</div>
    </div>
  )
}


// eslint-disable-next-line import/no-anonymous-default-export
export default {
  Resource,
  Project,
  Circle
}
