class MySocket{
  constructor(){
    this.mysocket =  null;
  }


  connectSocket(){
    console.log("memulai socket");
    var socket = new WebSocket("ws://localhost:8090/socket"); //make sure the port matches with your golang code
    this.mysocket = socket;

    socket.onmessage = (e)=>{  
      console.log("onmessage:",  e.data)
    }
    socket.onopen =  ()=> {
      console.log("socket opend")
      socket.send("greetings from js");
    }; 

    socket.onclose = ()=>{
      console.log("socket close")
    }
    
  }
}