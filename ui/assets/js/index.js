import HeaderIn from "./views/HeaderIn.js";
import HeaderOut from "./views/HeaderOut.js";
import Error from "./views/Error.js";
import Home from "./views/Home.js";
import LogIn from "./views/LogIn.js"
import SignUp from "./views/SignUp.js";
import Profile from "./views/Profile.js";
import ShowCat from "./views/ShowCat.js";
import PostView from "./views/PostView.js";
import NewPost from "./views/NewPost.js";
import Chat from "./views/Chat.js";
import connectSocket from "./views/Chat.js";



export {Router}

const pathToRegex = path => new RegExp("^" + path.replace(/\//g, "\\/").replace(/:\w+/g, "(.+)") + "$");
var authenticated = stringToBool(getCookie('auth'));
var authUserName = "undefined";
var wSocket = new MySocket()

var now = new Date();
var time = now.getTime();
var expireTime = time + 1000*36000;
now.setTime(expireTime);


const Router = async () => {

    const routes = [
        { path: "/", view: Home },
        { path: "/login", view: LogIn },
        { path: "/error", view: Error },
        { path: "/signup", view: SignUp },
        { path: "/profile", view: Profile },
        { path: "/category/:id", view: ShowCat },
        { path: "/logout", view: LogIn },
        { path: "/post/:id", view: PostView },
        { path: "/submit", view: NewPost },
    ];


    // Test each route for potential match
    const potentialMatches = routes.map(route => {
        return {
            route,
            result: location.pathname.match(pathToRegex(route.path))
        };
    });



    let match = potentialMatches.find(potentialMatch => potentialMatch.result !== null);

    /* Route not found - return first route OR a specific "not-found" route */
    if (!match) {
        match = {
            route: routes[2],
            result: [location.pathname]
        };
    }


    if (match.route.path == '/logout') {
        document.cookie = "auth=false;"
        // var now = new Date();
        // var time = now.getTime();
        // var expireTime = time + 1000 * 36000;
        // now.setTime(expireTime);

        fetch("/logout", {
            method: "POST"
        });
    };



    authenticated = stringToBool(getCookie('auth'));
    // console.log("USER AUTHENTICATED: ", authenticated)

    if(!authenticated && match.route.path != '/signup' && match.route.path != '/login'){
        location.assign('http://localhost:8090/login')
    }


    const v = getParams(match);
    var dataUrl = "/data/"

    if (v.id != undefined && v.url != undefined) {
        dataUrl += v.url + v.id;
    }

    
    if (v.url == "/profile") {
        dataUrl = "/data/profile"
    }

    if (v.url == "/comment") {
        dataUrl = "/data/comment"
    }

    var data = await fetchData(dataUrl);
    const view = new match.route.view(data);
    document.querySelector("#app").innerHTML = await view.getHtml();


debugger

    if (authenticated) {
        const chat = new Chat(data);
        const headin = new HeaderIn();
        
        document.querySelector("#header").innerHTML = await headin.getHtml();
        if (document.querySelector("#messageDiv").innerHTML == "") {

            authUserName = data.AuthenticatedUser != null? data.AuthenticatedUser.Name:  "unauthenticated";

            let message = {messageType:"online", context: authUserName};
            let msg = JSON.stringify(message);
            var socket = new MySocket()

            socket.connectSocket(msg);
            document.querySelector("#messageDiv").innerHTML = await chat.getHtml();
        }
    } else {
        const headout = new HeaderOut();
        document.querySelector("#messageDiv").innerHTML == ""
        document.querySelector("#header").innerHTML = await headout.getHtml();
    }

};


const navigateTo = url => {
    history.pushState(null, null, url);
    Router();
};


document.addEventListener("DOMContentLoaded", () => {
    document.body.addEventListener("click", e => {

        if (e.target.matches("[data-link]")) {

            e.preventDefault();
            navigateTo(e.target.href);
        }

    });

    Router();
});


window.addEventListener("popstate", Router);


window.addEventListener("beforeunload", function(){
    let message = {messageType:"offline", context: authUserName};
    let msg = JSON.stringify(message);
    wSocket.mysocket.send(msg);

})


const getParams = match => {
    const values = match.result.slice(1);
    const keys = Array.from(match.route.path.matchAll(/:(\w+)/g)).map(result => result[1]);

    var obj = Object.fromEntries(keys.map((key, i) => {
        return [key, values[i]];
    }));

    var url = match.route.path.match(/(?<=\/).+?(?=:)/g);

    if (url != null) {
        obj['url'] = url[0]
    } else {
        obj['url'] = match.route.path
    }


    return obj
};


async function fetchData(url) {

    var obj = fetch(url, {
        method: 'POST'
    })
        .then(response => {

            if (!response.ok) {
                throw new Error(`HTTP error: ${response.status}`);
            }
            // Otherwise (if the response succeeded), our handler fetches the response
            // as text by calling response.text(), and immediately returns the promise
            // returned by `response.text()`.
            return response.text()

        })
        .then(json => JSON.parse(json))
        .catch(err => console.error(`Fetch problem: ${err.message}`))

    return obj
}



function getCookie(cname) {
    let name = cname + "=";
    let decodedCookie = decodeURIComponent(document.cookie);
    let ca = decodedCookie.split(';');
    for(let i = 0; i <ca.length; i++) {
      let c = ca[i];
      while (c.charAt(0) == ' ') {
        c = c.substring(1);
      }
      if (c.indexOf(name) == 0) {
        return c.substring(name.length, c.length);
      }
    }
    return "";
  }


function stringToBool(str){
    if  (str == 'true'){
        return true
    }
        return false

}