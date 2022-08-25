import React, { useEffect, useState } from 'react';
import axios from 'axios';
import {useNavigate} from 'react-router-dom';
import Cookies from 'universal-cookie';

//similar to SignupPage
const LoginPage = () => {
    const navigate = useNavigate();
    const cookies = new Cookies();

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
                if (response.data.msg === "Login successful. Redirecting to home page") {
                    //set cookie loggedin to true
                    cookies.set('loggedin', '1')
                    //create cookie currentuser with username as value
                    cookies.set('currentuser', username);

                    setTimeout(() => {
                        navigate('/')
    
                    },2000)
    
                }
            })
        }
    }

    const enterPressed = event => {
        if (event.key === 'Enter') {
            event.preventDefault();

            login();
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
                onKeyDown={enterPressed}

                />
            {/* Password field */}
                <label htmlFor='password'>Password</label>
                <input 
                type="password" 
                name="password" 
                id="password" 
                placeholder="Enter password" 
                onChange={passwordChange}
                onKeyDown={enterPressed}
                />
                <input type="button" value="Login" onClick={login} />
            </form>
            {/* Display status message */}
            <p>
                {status}
            </p>


        </div>
    )

}

export default LoginPage