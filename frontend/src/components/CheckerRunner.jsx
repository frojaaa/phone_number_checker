import React from "react"
import {Button, Form} from "react-bootstrap";
import * as yup from "yup";
import {useForm} from "react-hook-form";
import {yupResolver} from "@hookform/resolvers/yup";

const schema = yup.object({
    numWorkers: yup.number().positive().integer().max(100).required(),
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
            type: "text"
        },
        {
            name: "lkLogin",
            promptText: "Логин ЛК",
            type: "text"
        },
        {
            name: "lkPassword",
            promptText: "Пароль ЛК",
            type: "password"
        },
        {
            name: "botToken",
            promptText: "Токен бота",
            type: "password"
        },
        {
            name: "tgUserID",
            promptText: "ID пользователя в тг",
            type: "text"
        },
        {
            name: "inputFileDir",
            promptText: "Путь к папке с входными файлами (лучше скопировать в проводнике)",
            type: "text"
        },
        {
            name: "outputFileDir",
            promptText: "Путь к папке с выходными файлами (лучше скопировать в проводнике)",
            type: "text"
        }
    ];
    const onSubmit = (data) => {
        console.log('data', data);
    }
    return (
        <Form noValidate onSubmit={handleSubmit(onSubmit)} className={"login-form"}>
            {fields.map(((field, index) => {
                return (
                    <Form.Group className="mb-3" controlId={field.name} key={index}>
                        <Form.Label className={"w-100 text-start mb-0"}>{field.promptText}</Form.Label>
                        <Form.Control type={field.type} {...register(field.name)} required
                                      onChange={handleChange} isInvalid={errors[field.name]}/>
                        <Form.Control.Feedback type={"invalid"}>{errors[field.name]?.message}</Form.Control.Feedback>
                    </Form.Group>
                )
            }))}
            <Button type="submit">Запустить</Button>
        </Form>
    )
}

export default CheckerRunner