import React from "react";
import { Routes, Route } from 'react-router-dom';
import SignupPage from "./Components/SignupPage";
import Home from "./Components/Home"
import LoginPage from "./Components/LoginPage"

const RouteHandler = () => {
    return (
        <Routes>
            <Route path='/' element={<Home/>}></Route>
            <Route path='/signup' element={<SignupPage/>}></Route>
            <Route path="/login" element={<LoginPage/>}></Route>

        </Routes>
    )
}

export default RouteHandler