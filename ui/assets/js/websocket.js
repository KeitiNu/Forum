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

      if (json.ContextType == "chat"){
        notify(json.Sender, json.Message)
        
      }else if (json.ContextType == "typing"){
        typing(json.Sender)
      }

      // if (json.Message) {
      //   notify(json.Sender, json.Message)
      // }else if (json.Sender) {
      //   typing(json.Sender)
      // }
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
}