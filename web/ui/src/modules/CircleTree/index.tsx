import React, { useEffect, useRef, useState } from 'react'
import { useParams } from 'react-router-dom';
import { DagreReact, Rect, RecursivePartial, NodeOptions, EdgeOptions, Size, ReportSize, ValueCache } from 'dagre-reactjs'
import nodes from './nodes'
import { Alert } from 'reactstrap';
import { getCircle, getCircleTree } from '../../core/api/circle';
import Sidebar from './Sidebar'
import './style.css'


type ElementsState = {
  nodes: Array<RecursivePartial<NodeOptions>>;
  edges: Array<RecursivePartial<EdgeOptions>>;
};

const DEFAULT_NODE_CONFIG = {
  styles: {
    node: {
      padding: {
        top: 10,
        bottom: 10,
        left: 10,
        right: 10
      }
    },
    shape: {},
    label: {}
  }
};

const getEdges = (projectId: string, resources: any[]) => {
  let allEdges: any[] = []
  for (let i = 0; i < resources.length; i++) {

    if (resources[i]?.parents) {
      const edges = resources[i]?.parents?.map((parent: any) => ({
        from: `${parent?.kind}-${parent?.name}`,
        to: `${resources[i]?.ref?.kind}-${resources[i]?.ref?.name}`
      }))

      allEdges = [...allEdges, ...edges]
    } else {
      let newEdge = {
        from: projectId,
        to: `${resources[i]?.ref?.kind}-${resources[i]?.ref?.name}`
      }

      allEdges = [...allEdges, newEdge]
    }
  }

  return allEdges
}


const getResources = (resources: any[]) => {
  return resources?.map((resource, i) => {
    let newRes = {
      id: `${resource?.ref?.kind}-${resource?.ref?.name}`,
      label: `${resource?.ref?.kind}: ${resource?.ref?.name}`,
      styles: {
        shape: {
          styles: { fill: "#fff", stroke: "#000", strokeWidth: "0" }
        },
        node: {
          padding: {
            top: 0,
            bottom: 0,
            left: 0,
            right: 0
          }
        },
        label: {
          styles: { pointerEvents: "none" }
        },

      },
      labelType: "resource",
      meta: { ...resource?.ref }
    }

    if (resource?.ref?.health) {
      let health = resource?.ref?.health
      if (health.status !== "Healthy") {
        newRes.styles.shape.styles.stroke = "#ff0000"
      }
    }

    return newRes
  })
}

const getElements = (circleName: string, nodes: any[]): ElementsState => {
  var nod: any[] = []
  var edges: any[] = []
  var circleNode = {
    id: `circle-${circleName}`,
    label: `Circle ${circleName}`,
    styles: {
      shape: {
        styles: { fill: "#fff", stroke: "#000", strokeWidth: "0" }
      },
      node: {
        padding: {
          top: 0,
          bottom: 0,
          left: 0,
          right: 0
        }
      },
      label: {
        styles: { pointerEvents: "none" }
      }
    },
    labelType: "circle",
  }

  for (let i = 0; i < nodes.length; i++) {
    let projectNodeID = `project-${nodes[i]?.name}`
    let projectNode = {
      id: projectNodeID,
      label: `project ${nodes[i]?.name}`,
      styles: {
        shape: {
          styles: { fill: "#fff", stroke: "#000", strokeWidth: "0" }
        },
        node: {
          padding: {
            top: 0,
            bottom: 0,
            left: 0,
            right: 0
          }
        },
        label: {
          styles: { pointerEvents: "none" }
        }
      },
      labelType: "project",
    }

    let newProjectCircleEdge = {
      from: `circle-${circleName}`,
      to: projectNodeID,
    }

    nod = [...nod, ...getResources(nodes[i]?.resources), projectNode, circleNode]
    edges = [...edges, ...getEdges(projectNodeID, nodes[i]?.resources), newProjectCircleEdge]
  }

  return {
    nodes: nod,
    edges: edges,
  }
}


const CircleTree = () => {
  const containerRef = useRef<any>(null)
  const { name } = useParams<any>()
  const [circle, setCircle] = useState<any>({})
  const [elements, setElements] = useState<ElementsState>({ nodes: [], edges: [] })
  const [stage, setStage] = useState(0)
  const [exception, setException] = useState("")

  const getCircleTreeReq = async () => {
    try {
      const circleTreeRes = await getCircleTree(name)
      setElements(getElements(circle?.name, circleTreeRes?.nodes))
      setStage(stage => stage + 1)
      setTimeout(getCircleTreeReq, 3000)
    } catch (e) {
      setException('Cannot get circle resource tree: ' + e.message)
    }

  }

  const getCircleReq = async () => {
    try {
      const circleRes = await getCircle(name)
      await setCircle(circleRes)
    } catch (e) {
      alert("Cannot get circle data: " + e.message)
    }
  }

  useEffect(() => {
    (async () => {
      await getCircleReq()
      getCircleTreeReq()
    })()
  }, [])

  return (
    <>
      <Sidebar
        circle={circle}
      />
      <div className="circle-tree-container" ref={containerRef}>
        {exception !== "" && (
          <Alert color="danger">
            {exception}
          </Alert>
        )}
        <svg id="schedule" width="100%" height="100%">
          <DagreReact
            nodes={elements.nodes}
            edges={elements.edges}
            defaultNodeConfig={DEFAULT_NODE_CONFIG}
            customNodeLabels={{
              "resource": {
                renderer: nodes.Resource,
                html: true
              },
              "project": {
                renderer: nodes.Project,
                html: true
              },
              "circle": {
                renderer: nodes.Circle,
                html: true
              }
            }}
            graphOptions={{
              marginx: 100,
              marginy: 15,
              rankdir: "LR",
              ranksep: 55,
              nodesep: 15
            }}
            stage={stage}
          />
        </svg>
      </div>
    </>
  )
}

export default CircleTree