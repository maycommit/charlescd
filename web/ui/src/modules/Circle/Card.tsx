import { pointer } from 'd3'
import React from 'react'
import './style.scss'

const Card = ({ circle, onCircleClick }: any) => {

  return (
    <div className="circle__card">
      <div className="circle__card__header">
        <div onClick={() => onCircleClick(circle?.name)} style={{ cursor: 'pointer' }}>
          <strong>{circle?.name}</strong>
        </div>
        <small className={`circle__card__header__${circle?.status?.status}`}>{circle?.status?.status}</small>
      </div>
      <div className="circle__card__clusters">
        {circle?.status?.projects?.sort((a: any, b: any) => a.name.localeCompare(b.name))
          .map((project: any) => (
            <div className={`circle__card__clusters__item__${project?.status}`}>
              {project?.name}
            </div>
          ))}
      </div>
    </div>
  )
}

export default Card