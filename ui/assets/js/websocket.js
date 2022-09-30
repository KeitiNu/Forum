class MySocket{
  constructor(){
    this.mysocket =  null;
    // this.counter = 0;
  }


  connectSocket(msg){
    console.log("memulai socket");
    console.log(msg);

    var socket = new WebSocket("ws://localhost:8090/socket"); //make sure the port matches with your golang code
    this.mysocket = socket;

    socket.onmessage = (e)=>{  
      console.log("onmessage:",  e.data)
    }

    socket.onopen =  ()=> {
      console.log("socket opend")

      socket.send(msg);
      // this.counter++;
    }; 

    socket.onclose = ()=>{
      console.log("socket close")
    }
    
  }


  // sendMessage(msg){
  //   // var string = msg.toString()
  //   console.log(msg)
  //   this.mysocket.send("send messagde");
  // }

}