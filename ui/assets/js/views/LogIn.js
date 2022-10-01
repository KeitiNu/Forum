import AbstractView from "./AbstractView.js";

export default class extends AbstractView {
    constructor(params) {
        super(params);
        this.setTitle("Kodify - Login");
    }


    get login(){
        $(document.body).on('submit', 'form#loginform', async function (e) {
            debugger
            e.preventDefault();
    
            var data = new FormData(e.target);
            var values = Object.fromEntries(data.entries());
    
            const location = window.location.pathname
            var o = await fetchFormData(values, location)
    
            this.params = o;

                const errors = this.params.Form.Errors.Errors
                const keys = Object.keys(errors)
    
                if (keys.length == 0) {
                    const tempLink = document.createElement('a')
                    const tempLocation = document.querySelector('.registerlink')

                    if (o.AuthenticatedUser != null) {
                        document.cookie = "auth=true;"

debugger
                        var mysocket = new MySocket()
                        mysocket.connectSocket(o.AuthenticatedUser.Name);
                        // var mysocket = new MySocket()
                        // mysocket.connectSocket();
                        // mysocket.sendMessage(o.AuthenticatedUser)
                        // mysocket.send(o.AuthenticatedUser);

                        tempLink.href = '/'
                        tempLink.dataset.link

                        tempLocation.appendChild(tempLink)
                        tempLink.click()
                    }else{
                        $('#errorgeneral').text("Unable to login")
                    }
    
                }else{
                    var errorSpots = document.querySelectorAll('.error')
    
                    errorSpots.forEach(err => {
                        err.innerHTML = ""
                    });
    
                    keys.map(function(key){
                        var spot = $('#error'+key)
                        spot.text(errors[key])
                    })
                }
    
        });
    
    
        async function fetchFormData(value, url) {
            var obj = fetch('/data'+url, {
                method: 'POST',
                headers: {
                    'Content-type': 'application/json; charset=UTF-8'
                },
                body: JSON.stringify(value)
            })
                .then(response => {

                    if (!response.ok) {
                        throw new Error(`HTTP error: ${response.status}`);
                    }
                    // Otherwise (if the response succeeded), our handler fetches the response
                    // as text by calling response.text(), and immediately returns the promise
                    // returned by `response.text()`.

                    // console.log("RESPONSEtext:", response.text())

                    return response.text()
    
                })
                .then(json => {
                    // console.log("RESPONSEtext:", json)
                    
                    
                    return JSON.parse(json)})
                .catch(err => console.error(`Fetch problem: ${err.message}`))
    

                
            var ans = obj.then(function (result) {
                return result; // "Some User token"
            })

            return obj
        }
    
    };

    async getHtml() {
        return `
        ${this.login}
        <div class="loginboxcontent">
        <div class=loginbox>
            <div class="iconbox">
                <span class="iconboxspan">
                    <img src="/static/css/images/logo_white.png">
                </span>
                <form method="post" id="loginform" class="loginform">
                    <div>
                        <label class="form-label" for="username"></label>
                        <input class="form-control" type="text" id="username" name="username" placeholder="Email or username" autocomplete="username"
                        <label class='error' id="errorusername"></label>
                    </div>
                    <div>
                            <label class="form-label" for="password"></label>
                            <input class="form-control" type="password" id="password" name="password" placeholder="Password" autocomplete="current-password">
                            <label class='error' id='errorpassword'></label>
                            
                            <label class='error' id='errorgeneric'></label>
                    </div>
                    <input class="btn loginbtn" type="submit" value="Log In">

                    <div class="registerlink">
                        <a href="/signup" data-link>New user? Register</a>
                    </div>
                    
                </form>
            </div>
        </div>
    </div>

    `;
    }
}