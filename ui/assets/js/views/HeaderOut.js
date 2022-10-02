import AbstractView from "./AbstractView.js";

export default class extends AbstractView {
	constructor(params) {
		super(params);
	}

	async getHtml() {
		return `

        <nav class="navbar header py-1 navbar-expand-sm fixed-top">
	<div class="container-fluid">
	  <a class="navbar-brand" style="width: 200px;" href="/" data-link>
		<img id="header-mainlogo" src="static/css/images/logo_white.png"  class="rounded-pill"></a>
			<form class="d-flex">
				<div class="btn-header0-div">
					<a class="btn header btn-outline-secondary first" href="/login" data-link>Log In	</a>
				</div>
				<div class="btn-header1-div">
					<a class="btn header btn-outline-secondary second" href="/signup" data-link>Sign Up</a>
				</div>
			</form>
	</div>
</nav>
        `;
	}
}