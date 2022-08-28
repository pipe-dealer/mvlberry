import { useLocation, useNavigate } from "react-router-dom"
import Cookies from 'universal-cookie';
import { closeconn, newconn, sendMsg } from "./webconn";
import { useEffect, useState } from "react";

const cookies = new Cookies();

const Chat = () => {
    const navigate = useNavigate();
    const [rec_msgs, setrMsg] = useState([]) //messages received from other user
    const [user_msg, setuMsg] = useState('') //message to be sent to other user
    const currentuser = cookies.get('currentuser')
    const loc = useLocation();

    const getmsg = (msg) => {
        setrMsg(rec_msgs => [...rec_msgs,msg])
    }

    //displays received messages
    const Displaymsg = () => {
        const rec_msg = rec_msgs.map((msg,index) => (
            <p key={index}>{msg}</p>
        ));

        return(
            <div>
                {rec_msg}
            </div>
        )
    }

    const msgChange = (event) => {
        let { value} = event.target
        setuMsg(value)
    }

    //gets message from user_msg and sends to other user
    const send = event => {
        if (event.key === 'Enter') {
            event.preventDefault()

            sendMsg(user_msg)
        }
    }

    //start websocket connection with endpoint
    useEffect(() => {
        newconn(currentuser,loc.state.user2,getmsg)
    }, [])

    //disconnect from chat and return to homepage
    const disconnect = () => {
        closeconn()
        navigate('/')
    }
    
    return(
        <div>
            <h1>Chatting with {loc.state.user2}</h1>
            <input type='button' value='disconnect' onClick={disconnect}></input>
            <input
            type = 'text'
            name = 'umessage'
            id = 'usermessage'
            placeholder="Press enter to send a message"
            onChange={msgChange}
            onKeyDown={send} />
            <br></br>
            <Displaymsg />
        </div>

    )
}

export default Chat