let msg = {
  name: 'channel add',
  data: {
    name: 'Hardware Support'

  }
}

let subMsg = {
  name: 'channel subscribe'
}

let ws = new WebSocket('ws://localhost:4000');

ws.onopen = () => {
  ws.send(JSON.stringify(subMsg));
  ws.send(JSON.stringify(msg));
}

ws.onmessage = (e) => {
  console.log(e.data);
}
