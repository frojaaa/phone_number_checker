import React from "react"
import {Button, Form} from "react-bootstrap";
import * as yup from "yup";
import {useForm} from "react-hook-form";
import {yupResolver} from "@hookform/resolvers/yup";
import CheckerService from "../checker.service"

const schema = yup.object({
    numWorkers: yup.number().positive().integer().max(101).required(),
    lkLogin: yup.string().required(),
    lkPassword: yup.string().required(),
    botToken: yup.string().required(),
    tgUserID: yup.number().required(),
    inputFileDir: yup.string().required(),
    outputFileDir: yup.string().required(),
}).required();

function CheckerRunner({userState, setUserState}) {
    const {register, handleSubmit, formState: {errors}} = useForm({
        resolver: yupResolver(schema)
    });
    const handleChange = e => {
        const {name, value} = e.target;
        setUserState(prevState => ({
            ...prevState,
            [name]: value
        }));
    };
    const fields = [
        {
            name: "numWorkers",
            promptText: "Количество потоков (<=100)",
            type: "text",
            errorText: "Значение должно быть целым и не более 100"
        },
        {
            name: "lkLogin",
            promptText: "Логин ЛК",
            type: "text",
            errorText: "Поле не должно быть пустым"
        },
        {
            name: "lkPassword",
            promptText: "Пароль ЛК",
            type: "password",
            errorText: "Поле не должно быть пустым"
        },
        {
            name: "botToken",
            promptText: "Токен бота",
            type: "password",
            errorText: "Поле не должно быть пустым"
        },
        {
            name: "tgUserID",
            promptText: "ID пользователя в тг",
            type: "text",
            errorText: "Поле не должно быть пустым и должно быть строковым"
        },
        {
            name: "inputFileDir",
            promptText: "Путь к папке с входными файлами (лучше скопировать в проводнике)",
            type: "text",
            errorText: "Поле не должно быть пустым и должно быть строковым"
        },
        {
            name: "outputFileDir",
            promptText: "Путь к папке с выходными файлами (лучше скопировать в проводнике)",
            type: "text",
            errorText: "Поле не должно быть пустым и должно быть числом"
        }
    ];
    const onSubmit = (data) => {
        console.log(errors)
        CheckerService.runChecker(
            data.numWorkers, data.lkLogin, data.lkPassword, data.botToken,
            data.tgUserID, data.inputFileDir, data.outputFileDir
        ).then(r => {
            if (r.status === 200) {
                setUserState(prevState => ({
                    ...prevState,
                    submitSuccess: true
                }));
            }
        })
    }
    return (
        <Form noValidate onSubmit={handleSubmit(onSubmit)} className={"login-form"}>
            {fields.map(((field, index) => {
                return (
                    <Form.Group className="mb-3" controlId={field.name} key={index}>
                        <Form.Label className={"w-100 text-start mb-0"}>{field.promptText}</Form.Label>
                        <Form.Control type={field.type} {...register(field.name)} required
                                      onChange={handleChange} isInvalid={errors[field.name]}/>
                        <Form.Control.Feedback type={"invalid"}>{errors[field.name] ? field.errorText : null}</Form.Control.Feedback>
                    </Form.Group>
                )
            }))}
            <Button type="submit" variant={userState.submitSuccess ? "success" : "outline-primary"}>Запустить</Button>
        </Form>
    )
}

export default CheckerRunner