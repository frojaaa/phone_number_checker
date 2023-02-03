import React, {useState} from "react";
import {Button, Form} from "react-bootstrap";
import AuthService from "../auth.service"

function Login({userState, setUserState}) {
    const [validated, setValidated] = useState(false);
    const handleSubmit = async (event) => {
        const form = event.currentTarget;
        event.preventDefault();
        if (form.checkValidity() === false) {
            event.stopPropagation();
        }

        setValidated(true);
        await AuthService.login(userState.username, userState.password);
        setUserState(prevState => ({
            ...prevState,
            password: "",
            enteredCheckerPassword: localStorage.getItem("enteredCheckerPassword")
        }));
        // window.location.reload();
    };
    const handleChange = e => {
        const {name, value} = e.target;
        setUserState(prevState => ({
            ...prevState,
            [name]: value
        }));
    };
    return (
        <>
            <Form noValidate validated={validated} onSubmit={handleSubmit} className={"login-form"}>
                <Form.Group className="mb-3" controlId="formGroupEmail">
                    <Form.Label className={"w-100 text-start mb-0"}>Логин</Form.Label>
                    <Form.Control type={"text"} placeholder="Введите логин" name={"username"} required
                                  onChange={handleChange}/>
                    <Form.Control.Feedback type={"invalid"}>Пожалуйста, введите имя пользователя</Form.Control.Feedback>
                </Form.Group>
                <Form.Group className="mb-3" controlId="formGroupPassword">
                    <Form.Label className={"w-100 text-start mb-0"}>Пароль</Form.Label>
                    <Form.Control type="password" name={"password"} placeholder="Введите пароль" required
                                  onChange={handleChange}/>
                    <Form.Control.Feedback type={"invalid"}>Пожалуйста, пароль</Form.Control.Feedback>
                </Form.Group>
                <Button type="submit">Войти</Button>
            </Form>
        </>
    )
}

export default Login