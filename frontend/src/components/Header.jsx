import React from "react";
import {Container, Nav, Navbar} from "react-bootstrap";
import {NavLink} from "react-router-dom";
import AuthService from "../auth.service"

function Header() {
    let token = localStorage.getItem("token")
    return (
        <header>
            <Navbar bg="dark" expand="lg">
                <Container>
                    <Navbar.Toggle aria-controls="basic-navbar-nav"/>
                    <Navbar.Collapse id="basic-navbar-nav">
                        <Nav className="ms-auto">
                            <NavLink className={"nav-link"} to="/">Главная</NavLink>
                            {
                                token === null ? null :
                                    <NavLink className={"nav-link"} to="#" onClick={async e => {
                                    await AuthService.logout();
                                    await window.location.reload();
                                    }
                                    }>Выйти</NavLink>
                            }
                        </Nav>
                    </Navbar.Collapse>
                </Container>
            </Navbar>
        </header>
    )
}

export default Header