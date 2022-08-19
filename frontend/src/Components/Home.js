import React, { useState } from 'react';
import axios from 'axios';
import {useNavigate} from 'react-router-dom';
import Cookies from 'universal-cookie';

const Home = () => {
    const cookies = new Cookies();
    //check if user has logged in
    if (cookies.get('loggedin') == 1) {
        //displays a welcome message with username
        return(
            <div>
                <h1>
                    Welcome {cookies.get('currentuser')}
                </h1>
            </div>
        )
    } else {
        return (
            <h1>
                Home
            </h1>
        )
    }
}

export default Home