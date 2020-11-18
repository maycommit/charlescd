import React, { useState } from 'react'
import { Link } from 'react-router-dom';
import {
  Collapse,
  Navbar as ReactstrapNavbar,
  NavbarToggler,
  NavbarBrand,
  Nav,
  NavItem,
  NavLink,
  UncontrolledDropdown,
  DropdownToggle,
  DropdownMenu,
  DropdownItem,
  NavbarText
} from 'reactstrap';

const Navbar = () => {
  const [isOpen, setIsOpen] = useState(false);

  const toggle = () => setIsOpen(!isOpen);

  return (
    <ReactstrapNavbar color="dark" dark expand="md">
      <NavbarBrand href="/">Charles CD</NavbarBrand>
      <NavbarToggler onClick={toggle} />
      <Collapse isOpen={isOpen} navbar>
        <Nav className="mr-auto" navbar>
          <NavItem>
            <NavLink><Link to="/projects">Projects</Link></NavLink>
          </NavItem>
          <NavItem>
            <NavLink><Link to="/circles">Circles</Link></NavLink>
          </NavItem>
        </Nav>
        <NavbarText>Simple Text</NavbarText>
      </Collapse>
    </ReactstrapNavbar>
  )
}

export default Navbar