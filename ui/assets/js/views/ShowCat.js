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
                    ${this.params.Categories.map(function(cat){
                        return "<h1 class='kodify'>"+cat.Title+"</h1><h2>"+cat.Description+"</h2>"
                    })}
                    </div>
                </div>
                <div class="adderbuttons">
                    <a class="btn adderbutton" href="/submit" data-link>Create Post</a>
                </div>
                <div class="categories">
                    <div class="insidecategories">

                    ${this.params.Posts.map(function(post){
                       var d = moment(post.Created).format("DD.MM.YYYY HH:mm"); 
                        return `  <div class="insidecatepadding">
                        <div class="catecard">
                            <div class="card-body">
                                <div class="post-card" id="post+`+post.ID+`">
                                    <div class="postdetails">
                                        <div class="post-username">Posted by `+post.User+` `+d+`</div>
                                        <div class="post-title"><a class="post-title stretched-link" href="/post/`+post.ID+`" data-link>`+post.Title+`</a></div>
                                        <div class="post-description">`+post.Content+`</div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>`
                    })}
                    </div>
                </div>
            </div>
        </div>
    </div>
        `;
    }
}

