let ws;

//gets responses from webserver
let connect = (ws,displaymsg) => {
    console.log('connecting to websocket endpoint...')

    ws.onopen = () => {
        console.log('connection established')
    }

    //received messages are passed through function in Chat.js
    ws.onmessage = msg => {
        displaymsg(msg.data)
    }

    ws.onerror = error => {
        console.log("Socket error: ", error)
    }
}

//create new websocket client
let newconn = (user1,user2,displaymsg) => {
    let url = encodeURI(`ws://localhost:8080/api/ws?user1=${user1}&user2=${user2}`)
    ws = new WebSocket(url)

    connect(ws,displaymsg)
}


let sendMsg = (msg) => {
    ws.send(msg)
}

let closeconn = () => {
    console.log('closing connection...')
    ws.close();
    console.log('connection closed')
}

export { newconn, sendMsg, closeconn }