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
                <form method="post" class="std auth loginform">
                    <div>
                        <label class="form-label" for="username"></label>
                        <input class="form-control" type="text" id="username" name="username" placeholder="Username"
                            value="">
                    </div>
                    <div>
                        <label class="form-label" for="email"></label>
                        <input class="form-control" type="text" id="email" name="email" placeholder="Email" value="">
                    </div>
                    <div>
                        <label class="form-label" for="pass"></label>
                        <input class="form-control" type="password" id="pass" name="password" placeholder="Password">
                    </div>
                    <div>
                        <label class="form-label" for="pass"></label>
                        <input class="form-control" type="password" id="password" placeholder="Confirm password"
                            name="confirm_password">
                    </div>
                    <input class="btn loginbtn" type="submit" value="Sign Up" data-link>
                    <div class="registerlink">
                        <a href="/login" data-link>Already registered? Log In</a>
                    </div>
                </form>
            </div>
        </div>
    </div>

    `;
    }
}