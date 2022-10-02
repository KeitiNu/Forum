import AbstractView from "./AbstractView.js";

export default class extends AbstractView {
    constructor(params) {
        super(params);
        this.setTitle("Kodify - Profile");
    }

    async getHtml() {
        return `

<div class="createpostcontent">
	<div class="profilebox">
		<div class="createpostboxinside">
			<div class="createpostheader">
				<h1>Profile</h1>
			</div>
			<div class="accounth2">
				<h2>Account details</h2>
			</div>

			<div class="accountdetails">
				<div class="profilename">
					<h3>Username</h3>
					<a>${this.params.AuthenticatedUser.Name}</a>
				</div>
				<div class="profileemail">
					<h3>Email</h3>
					<a>${this.params.AuthenticatedUser.Email}</a>
				</div>
			</div>
		</div>
			</div>
		</div>
	</div>
</div>
`;
}
}