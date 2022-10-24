import AbstractView from "./AbstractView.js";

export default class extends AbstractView {
    constructor(params) {
        super(params);
        this.setTitle("Kodify - Login");
    }

    async getHtml() {
        return `
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