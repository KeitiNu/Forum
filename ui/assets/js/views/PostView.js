import AbstractView from "./AbstractView.js";

export default class extends AbstractView {
    constructor(params) {
        super(params);
        this.setTitle("Kodify - Post");
    }

    async getHtml() {
        return `
        <div class="mainpagecontent">
        <div class="mainpagebox">
            <div class="mainpageboxinside">
                <div class="insidecateboxheader">
                    <div class="mainbpageboxheaderthread">
                        <div class="post-card" id="post${this.params.Post.ID}">
    
                            <div class="postdetails">
                                <div class="post-username-thread">Posted by ${this.params.Post.User} {{timeAgo .Created}}</div>
                                <div class="post-title-thread"><a class="post-title-thread stretched-link"
                                        href="/post/${this.params.Post.ID}" data-link>${this.params.Post.Title}</a></div>
                                <div class="post-description-thread">${this.params.Post.Content}</div>
                            </div>
                        </div>
                    </div>
                    <div class="categories">
                        <div class="commentingbox">
                            <form class="comment" method="POST">
                                <textarea class="commentbox " name="comment" id="" placeholder="Write your comment here"
                                    maxlength="2050"></textarea>
                                <div class="com">
                                    <button type="submit" class="btn submitbtn" name="submitPost">Submit comment</button>
                                </div>
                             
                        </div>
    
                        <div class="insidecategories">

                        ${this.params.Comments.map(function(comment){
                            var d = moment(comment.Created).format("DD.MM.YYYY HH:mm"); 

                            return ` 
                            <div class="insidecatepadding">
                                <div class="catecard">
                                    <div class="card-body">
                                        <div class="post-card" id="comment`+ comment.ID +`">
    
                                            <div class="postdetails">
                                                <div class="post-username">Posted by `+ comment.User +` `+ d +`</div>
                                                <p class="post-description-comment card-text">`+ comment.Content +`</p>
                                            </div>
                                  
                                        </div>
                                    </div>
                                </div>
                            </div>
                            <div class="modal fade" id="exampleModal" tabindex="-1" role="dialog"
                                aria-labelledby="exampleModalLabel" aria-hidden="true">
                                <div class="modal-dialog" role="document">
                                <div class="modal-content">
                                    <div class="modal-header">
                                        <h5 class="modal-title" id="exampleModalLabel">Edit comment</h5>
                                        <button type="button" class="close" data-dismiss="modal"
                                            aria-label="Close">
                                            <span aria-hidden="true">&times;</span>
                                        </button>
                                    </div>
                                    <div class="modal-body">
                                        <form method="POST">
                                            <label for="comment">Comment</label>
                                            <textarea name="commentUpdate" id="" cols="30"
                                                rows="10">{{.Content}}</textarea>
                                            <textarea name="commentUpdateID" id="" cols="0" rows="0"
                                                hidden>{{.ID}}</textarea>
                                            <textarea name="commentUpdateUser" id="" cols="0" rows="0"
                                                hidden>{{.User}}</textarea>
                                            <div class="modal-footer">
                                                <button type="button" class="btn btn-secondary"
                                                    data-dismiss="modal">Close</button>
                                                <button type="submit" class="btn btn-primary"
                                                    name="submitPost">Save changes</button>
                                            </div>
                                        </form>
                                    </div>
    
                                </div>
                            </div>
                        </div>	`						
                    })}
                        
                        </div>
                    </div>
    
                </div>
    
    
            </div>
        </div>
    </div>
        `;
    }
}