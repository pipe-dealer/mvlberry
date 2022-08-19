import React, { useState } from 'react';
import axios from 'axios';
import {useNavigate} from 'react-router-dom';

//similar to SignupPage
const LoginPage = () => {
    const navigate = useNavigate();

    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [status, setStatus] = useState('');

    const usernameChange = (event) => {
        let {_, value} = event.target
        value = value.replace(/\s+/g, '');
        setUsername(value)
    }

    const passwordChange = (event) => {
        let {_, value} = event.target
        value = value.replace(/\s+/g, '');

        setPassword(value)
    }

    const login = () => {
        if (username.length === 0){
            setStatus("No username entered")
        } else if (password.length === 0) {
            setStatus("No password entered")
        // if fields are not empty, make a post request to our web server and returns a response message
        } else {
            axios.post('http://localhost:8080/api/login', {
                username: username,
                password: password
            }).then((response) => {
                setStatus(response.data.msg)
                if (response.data.msg == "Login successful. Redirecting to home page") {
                    setTimeout(() => {
                        navigate('/')
    
                    },2000)
    
                }
            })
        }
    }

    return(
        <div>
            <form>
            {/* Username field */}
                <label htmlFor='username'>Username</label>
                <input 
                type="text" 
                name="username" 
                id="username" 
                placeholder="Enter username" 
                onChange={usernameChange}
                />
            {/* Password field */}
                <label htmlFor='password'>Password</label>
                <input 
                type="password" 
                name="password" 
                id="password" 
                placeholder="Enter password" 
                onChange={passwordChange}
                />
                <input type="button" value="Login" onClick={login}/>
            </form>
            {/* Display status message */}
            <p>
                {status}
            </p>


        </div>
    )

}

export default LoginPage