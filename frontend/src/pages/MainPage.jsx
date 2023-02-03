import React, {useState} from "react";
import {Container} from "react-bootstrap";
import Login from "../components/Login";
import CheckerPassword from "../components/CheckerPassword";
import Header from "../components/Header";
import CheckerRunner from "../components/CheckerRunner";

export default function MainPage() {
    const [userData, setUserData] = useState({username: "", password: ""})
    let token = localStorage.getItem("token")
    let enteredCheckerPassword = localStorage.getItem("enteredCheckerPassword") === "true"
    return (
        <>
            <Header/>
            <Container fluid style={{display: "grid", placeItems: "center"}}>
                {token === null ? <Login userState={userData} setUserState={setUserData}/> :
                    <CheckerRunner setUserState={setUserData}/>}
                {enteredCheckerPassword || token === null ? null :
                    <CheckerPassword userState={userData} setUserState={setUserData}/>}
            </Container></>
    )
}