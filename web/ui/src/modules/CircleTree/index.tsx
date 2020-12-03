import React, { useCallback, useEffect, useRef, useState } from 'react'
import { useParams } from 'react-router-dom';
import { DagreReact, RecursivePartial, NodeOptions, EdgeOptions } from 'dagre-reactjs'
import { UncontrolledReactSVGPanZoom } from 'react-svg-pan-zoom';
import AutoSizer from 'react-virtualized-auto-sizer';
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
        top: 0,
        bottom: 0,
        left: 0,
        right: 0
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
    let newRes: NodeOptions = {
      id: `${resource?.ref?.kind}-${resource?.ref?.name}`,
      label: `${resource?.ref?.kind}: ${resource?.ref?.name}`,
      shape: "rect",
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
        if (newRes?.styles?.shape?.styles) {
          newRes.styles.shape.styles.stroke = "#ff0000"
        }
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

  nod = [circleNode]

  for (let i = 0; i < nodes.length; i++) {
    let projectNodeID = `project-${nodes[i]?.name}`
    let projectNode = {
      id: projectNodeID,
      label: `Project ${nodes[i]?.name}`,
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

    nod = [...nod, ...getResources(nodes[i]?.resources), projectNode]
    nod = nod.sort((a, b) => a.label.localeCompare(b.label))
    edges = [...edges, ...getEdges(projectNodeID, nodes[i]?.resources), newProjectCircleEdge]
    edges = edges.sort((a, b) => a.to.localeCompare(b.to))
  }

  return {
    nodes: nod,
    edges: edges,
  }
}


const CircleTree = () => {
  const containerRef = useRef<any>()
  const viewer = useRef<any>()
  const { name } = useParams<any>()
  const [circle, setCircle] = useState<any>({})
  const [elements, setElements] = useState<ElementsState>({ nodes: [], edges: [] })
  const [stage, setStage] = useState(0)
  const [exception, setException] = useState("")
  const [size, setSize] = useState<any>({ width: 0, height: 0 })

  const getCircleTreeReq = useCallback(
    async () => {
      try {
        const circleTreeRes = await getCircleTree(name)
        const newElements = getElements(name, circleTreeRes?.nodes)
        setElements(newElements)
        setStage(stage => stage + 1)
        setTimeout(getCircleTreeReq, 2000)
      } catch (e) {
        setException('Cannot get circle resource tree: ' + e.message)
      }
    },
    [circle, name],
  )

  const getCircleReq = useCallback(
    async () => {
      try {
        const circleRes = await getCircle(name)
        setCircle(circleRes)
      } catch (e) {
        alert("Cannot get circle data: " + e.message)
      }
    },
    [name]
  )

  useEffect(() => {
    if (circle) {
      getCircleTreeReq()
    }
  }, [circle, getCircleTreeReq])

  useEffect(() => {
    if (circle) {
      return
    }
    getCircleReq()
  }, [getCircleReq, circle])


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
        <div style={{ height: "100%" }}>
          <AutoSizer>
            {({ height, width }: any) => (
              <UncontrolledReactSVGPanZoom
                width={width}
                height={height}
                tool="pan"
                background="#fff"
                detectAutoPan={false}
                miniatureProps={{
                  position: 'none',
                  background: '#fff',
                  width: 100,
                  height: 100,
                }}
                toolbarProps={{
                  position: 'none',
                  SVGAlignX: undefined,
                  SVGAlignY: undefined,
                }}
                ref={viewer}
              >
                <svg id="schedule" width={width} height={height}>
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
                      marginx: 150,
                      marginy: 15,
                      rankdir: "LR",
                      ranksep: 55,
                      nodesep: 15
                    }}
                    stage={stage}
                  />
                </svg>
              </UncontrolledReactSVGPanZoom>
            )}
          </AutoSizer>
        </div>
      </div >
    </>
  )
}

export default CircleTree