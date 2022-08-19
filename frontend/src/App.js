import React, {useState} from 'react';
import RouteHandler from './Routes';
import Cookies from 'universal-cookie';



const App = () => {
    const cookies = new Cookies();
    //set cookies to empty
    cookies.set('loggedin', 0);
    cookies.set('currentuser', '');

    return (
        <RouteHandler />
    )
}


export default App