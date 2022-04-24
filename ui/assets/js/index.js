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


const pathToRegex = path => new RegExp("^" + path.replace(/\//g, "\\/").replace(/:\w+/g, "(.+)") + "$");
var authenticated = true;



const router = async () => {


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
        authenticated = false;
    };


    if (match.route.path == '/login') {
        authenticated = true;
    };




    const v = getParams(match);
    var dataUrl = "/data/"

    if (v.id != undefined && v.url != undefined) {

        dataUrl += v.url + v.id ;
    }
    var data = await fetchData(dataUrl);
console.log(data)

    // if (match.route.path == '/') {
    const view = new match.route.view(data);
    document.querySelector("#app").innerHTML = await view.getHtml();
        
    // }else{

    //  const view = new match.route.view(getParams(match));

    // document.querySelector("#app").innerHTML = await view.getHtml();

    // }



    // if (match.route.path == '/') {
    // view.logData()

    // };

    if (authenticated) {

        const headin = new HeaderIn();
        document.querySelector("#header").innerHTML = await headin.getHtml();

    } else {
        
        const headout = new HeaderOut();
        document.querySelector("#header").innerHTML = await headout.getHtml();

    }

    // document.querySelector("#app").innerHTML = await view.getHtml();

};


const navigateTo = url => {
    history.pushState(null, null, url);
    router();
};


document.addEventListener("DOMContentLoaded", () => {
    document.body.addEventListener("click", e => {

        if (e.target.matches("[data-link]")) {

            e.preventDefault();
            navigateTo(e.target.href);
        }
    });

    /* Document has loaded -  run the router! */
    router();
});


window.addEventListener("popstate", router);


const getParams = match => {
    const values = match.result.slice(1);
    const keys = Array.from(match.route.path.matchAll(/:(\w+)/g)).map(result => result[1]);

    var obj =  Object.fromEntries(keys.map((key, i) => {
        return [key, values[i]];
    }));

    var url = match.route.path.match(/(?<=\/).+?(?=:)/g);
console.log(url)

    if (url != null) {
    obj['url'] = url[0]       
    }

    return obj
};



async function fetchData(url){
    console.log(url)
    var obj =  fetch(url, {
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

