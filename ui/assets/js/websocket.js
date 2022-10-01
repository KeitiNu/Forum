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
      var json = JSON.parse(e.data);
      var onlineUsers = json.OnlineUsers;


      console.log("onmessage:",  e.data)

      onlineUsers.forEach(user => {
        changeStatus(user);
      });

    }


    socket.onopen =  ()=> {
      console.log("socket opend")
      console.log("MSG:",  msg)


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