import dagreD3 from 'dagre-d3'
import * as d3 from 'd3'

export const getColorByHealth = (health: any): any => {
  if (health && health?.status) {
    switch (health?.status) {
      case 'Healthy':
        return 'node-success'
      case 'Degraded':
        return 'node-danger'
      case 'Progressing':
        return 'node-progress'
      default:
        return ''
    }
  }

  return ''
}

export const getTree = (data: any, circleId: string) => {
  let currentNodes: any = []
  let currentEdges: any = []


  let html = `<div class="node node-circle">`;
  html += `<div class=name>${circleId}</div>`;
  html += `<div class=kind>Circle</span></div>`;
  html += "</div>";
  currentNodes = [...currentNodes, {
    name: circleId,
    labelType: 'html',
    label: html,
    padding: 0,
    shape: 'rect',
    meta: {
      kind: 'Circle'
    }
  }]

  data?.projects?.sort((a: any, b: any) => a.name?.localeCompare(b.name)).map((project: any) => {
    var html = `<div class="node node-project" >`;
    html += `<div class=name>${project?.name}</div>`;
    html += `<div class=kind>Project</span></div>`;
    html += "</div>";
    currentNodes = [...currentNodes, {
      name: project?.name,
      labelType: 'html',
      label: html,
      padding: 0,
      shape: 'rect',
      meta: {
        kind: 'Project'
      }
    }]

    currentEdges = [...currentEdges, {
      source: circleId,
      target: project?.name,
      arrowhead: "normal"
    }]


    project?.resources?.map((res: any) => {
      let html = `<div class="node ${getColorByHealth(res?.ref?.health)}">`
      html += `<div class=name>${res?.ref?.name}</siv>`;
      html += `<div class=kind>${res?.ref?.kind}</span></div>`;
      html += '</div>'

      currentNodes = [...currentNodes, {
        name: `${res?.ref?.kind}-${res?.ref?.name}`,
        labelType: 'html',
        label: html,
        padding: 0,
        shape: 'rect',
        meta: {
          name: res?.ref?.name,
          kind: res?.ref?.kind,
          version: res?.ref?.version,
          group: res?.ref?.group,
          health: res?.ref?.health,
          creationTime: res?.ref?.creationTimestamp,
        }
      }]


      if (!res?.parents) {
        currentEdges = [...currentEdges, {
          source: project?.name,
          target: `${res?.ref?.kind}-${res?.ref?.name}`,
          arrowhead: "normal"
        }]

      } else {
        res?.parents?.map((parent: any) => {
          currentEdges = [...currentEdges, {
            source: `${parent?.kind}-${parent?.name}`,
            target: `${res?.ref?.kind}-${res?.ref?.name}`,
            arrowhead: "normal"
          }]
        })
      }
    })
  })

  return [currentNodes, currentEdges]
}
