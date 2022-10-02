class MySocket{
  constructor(){
    this.mysocket =  new WebSocket("ws://localhost:8090/socket");
    // this.counter = 0;
  }
  


  connectSocket(msg){
    console.log("memulai socket");
    console.log(msg);

    // var socket = new WebSocket("ws://localhost:8090/socket"); //make sure the port matches with your golang code
    // this.mysocket = socket;

    var socket = this.mysocket;

    socket.onmessage = (e)=>{  
      var json = JSON.parse(e.data);
      var onlineUsers = json.OnlineUsers ? json.OnlineUsers : [];


      console.log("onmessage:",  e.data)

      onlineUsers.forEach(user => {
        changeStatus(user, 1);
      });
      if (json.OfflineUser) {
        changeStatus(json.OfflineUser, 0);
      }
    }


    socket.onopen =  ()=> {
      console.log("socket opend")
      console.log("MSG:",  msg)

      socket.send(msg);
      // this.counter++;
    }; 


    socket.onclose = (e)=>{
      e.preventDefault();
      // socket.send({"Closing"});

      console.log("socket close")
    }
    
  }


  // sendMessage(msg){
  //   // var string = msg.toString()
  //   console.log(msg)
  //   this.mysocket.send("send messagde");
  // }

}