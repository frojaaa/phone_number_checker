import React from 'react';
import { Route, Routes } from "react-router-dom";
import {publicRoutes} from "../routes";
import {Navigate} from "react-router";


export default function () {
    return (
        <Routes>
            {publicRoutes.map(({ path, Component }) => {
                return <Route key={path} path={path} element={<Component/>} exact/>
            })}
            <Route path={"*"} element={<Navigate to={"/"} replace/>}/>
        </Routes>
    )
}

