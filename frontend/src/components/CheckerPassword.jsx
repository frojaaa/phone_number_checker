import React, {useState} from "react";
import {Button, Form} from "react-bootstrap";
import CheckerService from "../checker.service"

function CheckerPassword({userState, setUserState}) {
    const [validated, setValidated] = useState(false);
    const handleSubmit = async (event) => {
        const form = event.currentTarget;
        event.preventDefault();
        if (form.checkValidity() === false) {
            event.stopPropagation();
        }
        let username = JSON.parse(localStorage.getItem("username"))
        setValidated(true);
        await CheckerService.checkPassword(username, userState.password);
        localStorage.setItem("enteredCheckerPassword", "true")
        setUserState(prevState => ({
            ...prevState,
            enteredCheckerPassword: localStorage.getItem("enteredCheckerPassword")}));
    };
    const handleChange = e => {
        const {name, value} = e.target;
        setUserState(prevState => ({
            ...prevState,
            [name]: value
        }));
    };
    return (
        <Form noValidate validated={validated} onSubmit={handleSubmit} className={"login-form"}>
            <Form.Group className="mb-3" controlId="formGroupPassword">
                <Form.Label className={"w-100 text-start mb-0"}>Пароль чекера</Form.Label>
                <Form.Control type="password" name={"password"} placeholder="Введите пароль чекера" required
                              onChange={handleChange}/>
                <Form.Control.Feedback type={"invalid"}>Пожалуйста, пароль</Form.Control.Feedback>
            </Form.Group>
            <Button type="submit">Войти</Button>
        </Form>
    )
}

export default CheckerPassword;