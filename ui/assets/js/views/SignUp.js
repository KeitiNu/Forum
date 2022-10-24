import AbstractView from "./AbstractView.js";

export default class extends AbstractView {
    constructor(params) {
        super(params);
        this.setTitle("Kodify - Sign Up");
    }

   
    async getHtml() {
        return `
        <div class="loginboxcontent">
        <div class=loginbox>
            <div class="iconbox">
                <span class="iconboxspan">
                    <img src="/static/css/images/logo_white.png">
                </span>
                <form method="post" id="signupform" class="std auth loginform">
                    <div>
                        <label class="form-label" for="username"></label>
                        <input class="form-control" type="text" id="username" name="username" placeholder="Username" value="">
                        <label class='error' id="errorusername"></label>

                    </div>
            
                    <div>
                        <label class="form-label" for="email"></label>
                        <input class="form-control" type="text" id="email" name="email"  placeholder="Email" value="">
                        <label class='error' id="erroremail"></label>

                    </div>


                    <div>
					<label class="form-label" for="forname"></label>
					<input class="form-control" type="text" id="forname" name="forname" placeholder="First name">
                    <label class='error' id="errorforname"></label>

				</div>
				<div>
					<label class="form-label" for="surname"></label>
					<input class="form-control" type="text" id="surname" name="surname" placeholder="Last name">
                    <label class='error' id="errorsurname"></label>

				</div>
				<div class="form-outline">
					<label class="form-label white"></label>
					<input class="form-control" type="number" placeholder="Age" min='0' name='age'></input>
                    <label class='error' id="errorage"></label>

				</div>
				<div>
					<label class="form-label" for="gender"></label>
					<select name='gender' class="form-select form-select-sm select">
						<option selected hidden value="4">Gender</option>
						<option value="0">Female</option>
						<option value="1">Male</option>
						<option value="2">Non-binary</option>
						<option value="3">I prefer not to say</option>
					</select>
                    <label class='error' id="errorgender"></label>
				</div>

            
                    <div>
                        <label class="form-label" for="pass"></label>
                        <input class="form-control" type="password" id="pass" name="password" placeholder="Password">
                        <label class='error' id='errorpassword'></label>
                    </div>
                    <div>
                        <label class="form-label" for="pass"></label>
                        <input class="form-control" type="password" id="password" placeholder="Confirm password" name="confirmPassword">
                        <label class='error' id='errorconfirmPassword'></label>
                            
                        <label class='error' id='errorgeneric'></label>
                    </div>
            
                    <input class="btn loginbtn" type="submit" value="Sign Up">
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














