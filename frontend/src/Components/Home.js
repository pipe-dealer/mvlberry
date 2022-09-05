import React, { useEffect, useState } from 'react';
import axios from 'axios';
import {useNavigate} from 'react-router-dom';
import Cookies from 'universal-cookie';

const cookies = new Cookies();

//gets user's friends and displays them as buttons
const Displayfriends = () => {
    const [friendlist, setFriend] = useState([]) ;

    useEffect(() => {
        if (cookies.get('loggedin') == 1) {
            //get user's friends
            axios.get('http://localhost:8080/api/home/getfriends', {params: {current_user: cookies.get('currentuser')}}).then((response) => {
                const friends = []
                for (const r in response.data){
                    const friend = (response.data[r])
                    friends.push(friend)
                }
                setFriend(friends)
            })
        }
    }, []);

    return friendlist
}

const Home = () => {
    const navigate = useNavigate();

    //login button
    const loginbtn = () => {
        navigate('/login')
    }

    //sign up button
    const signupbtn = () => {
        navigate('/signup')
    }

    //check if user has logged in
    if (cookies.get('loggedin') == 1) {
        //displays a welcome message with username
        const friends = Displayfriends()
        const currentuser = cookies.get('currentuser')

        const startChat = (user2) => {
            navigate("/chat", {state : {user2: user2}})
        }

        return(
            <div>
                <h1>
                    Welcome {currentuser}
                </h1>
                <h2>
                    Your friends:
                </h2>
                <br></br>
                {friends.map((friend,index) => (
            <input key={index} type='button' value={friend} onClick={() => startChat(friend)}></input>
        ))}
            </div>
        )
    } else {
        return (
            <h1>
                Home
                <input type='button' value="Login" onClick={loginbtn} />
                <input type='button' value="Sign up" onClick={signupbtn} />
            </h1>
        )
    }
}

export default Home