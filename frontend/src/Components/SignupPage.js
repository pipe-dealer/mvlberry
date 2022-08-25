import React, { useState } from 'react';
import axios from 'axios';
import {useNavigate} from 'react-router-dom';

const SignupPage = () => {
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

    const createAccount = () => {

        // checks if username or password field is empty and returns an error message if it does
        if (username.length === 0){
            setStatus("No username entered")
        } else if (password.length === 0) {
            setStatus("No password entered")
        // if fields are not empty, make a post request to our web server and returns a response message
        } else {
            axios.post('http://localhost:8080/api/signup', {
                username: username,
                password: password
            }).then((response) => {
                setStatus(response.data.msg)
                if (response.data.msg == "Account successfully created. Redirecting to login page") {
                    setTimeout(() => {
                        navigate('/login')
    
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
                    placeholder="Enter desired username" 
                    onChange={usernameChange}
                    onKeyDown={enterPressed}

                    />
                {/* Password field */}
                    <label htmlFor='password'>Password</label>
                    <input 
                    type="password" 
                    name="password" 
                    id="password" 
                    placeholder="Enter a password" 
                    onChange={passwordChange}
                    onKeyDown={enterPressed}

                    />
                    <input type="button" value="Sign Up" onClick={createAccount}/>
                </form>
                {/* Display status message */}
                <p>
                    {status}
                </p>

            </div>

    )
    
}
export default SignupPage