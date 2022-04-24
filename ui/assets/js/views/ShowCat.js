import AbstractView from "./AbstractView.js";

export default class extends AbstractView {
    constructor(params) {
        super(params);
        this.setTitle("Kodify - Show Category");
        this.catId = params.id;

    }

    
    async getHtml() {
        return `
        <div class="mainpagecontent">
        <div class="mainpagebox">
            <div class="mainpageboxinside">
                <div class="insidecateboxheader">
                    <div class="headercard">
                        {{range .Categories}}
                            <h1 class="kodify">{{.Title}}</h1>
                            <h2>{{.Description}}</h2>
                        {{end}}
                    </div>
                </div>
                <div class="adderbuttons">
                    <a class="btn adderbutton" href="/" data-link>Back</a>
                    <a class="btn adderbutton" href="/submit" data-link>Submit Post</a>
                </div>
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
                    </div>
                </div>
            </div>
        </div>
    </div>
        `;
    }
}
