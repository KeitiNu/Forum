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
					<a>{{.AuthenticatedUser.Name}}</a>
				</div>
				<div class="profileemail">
					<h3>Email</h3>
					<a>{{.AuthenticatedUser.Email}}</a>
				</div>		
			</div>


			<div class="activity">
				<div class="accounth2">
					<h2>Your activity</h2>
				</div>

				<div class="">
					<a href="/profile/posts"><button type="button" class="btn submitbtn" data-link>Posts</button></a>
					<a href="/profile/comments"><button type="button" class="btn submitbtn" data-link>Comments</button></a>
				</div>

				<div class="mainboxinside">
					<div class="categories">
						<div class="insidecategories">
							{{range .Posts}}
							<div class="insidecatepadding">
								<div class="catecard">
									<div class="card-body">
										<div class="post-card" id="post{{.ID}}">
											<div class="postdetails">
												<div class="post-username">Posted by {{.User}} {{timeAgo .Created}}</div>
												<div class="post-title"><a class="post-title stretched-link" href="/post/{{.ID}}" data-link>{{.Title}}</a></div>
												<div class="post-description">{{.Content}}</div>
											</div>
										</div>
									</div>
								</div>
							</div>	
							{{end}}

							{{range .Comments}}
							<div class="insidecatepadding">
								<div class="catecard">
									<div class="card-body">
										<div class="postcard">
											<div class="votingbuttons">
											</div>
											<div class="postdetails">
												<div class="post-username">Posted by {{.User}} {{timeAgo .Created}}</div>
												<div class="post-description-thread">{{.Content}}</div>
											</div>
										</div>
									</div>
								</div>
							</div>
							{{end}}
						</div>
					</div>
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