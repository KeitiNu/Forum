import AbstractView from "./AbstractView.js";

export default class extends AbstractView {
    constructor(params) {
        super(params);
        this.postId = params.id;
        this.setTitle("Kodify - Post");
    }

    async getHtml() {
        return `
        <div class="mainpagecontent">
        <div class="mainpagebox">
            <div class="mainpageboxinside">
                <div class="insidecateboxheader">
                    <div class="mainbpageboxheaderthread">
                        <div class="post-card" id="post{{.ID}}">
                            <div class="votingbuttons">
                                <div class="post-upvote">
                                    <form id="form-up-{{.ID}}" onsubmit="return fetchpost('{{.ID}}', 'up')">
                                        <input value="{{.ID}}" hidden>
                                        <button type="submit" class="btn up-button"
                                            onclick="vote('post{{.ID}}', 'up-button')">
                                            <i class="bi upb bi-caret-up"></i>
                                        </button>
                                    </form>
                                </div>
    
                                <div class="votecounter">
                                    <a name="votes" class="votes" value="0">{{.Votes}}</a>
                                </div>
    
                                <div class="post-downvote">
                                    <form id="form-down-{{.ID}}" onsubmit="return fetchpost('{{.ID}}', 'down')">
                                        <input value="{{.ID}}" hidden>
                                        <button type="submit" class="btn down-button"
                                            onclick="vote('post{{.ID}}', 'down-button')">
                                            <i class="bi downb bi-caret-down"></i>
                                        </button>
                                    </form>
                                </div>
                            </div>
    
                            <div class="postdetails">
                                <div class="post-username-thread">Posted by {{.User}} {{timeAgo .Created}}</div>
                                <div class="post-title-thread"><a class="post-title-thread stretched-link"
                                        href="/post/{{.ID}}" >{{.Title}}</a></div>
                                <div class="post-description-thread">{{.Content}}</div>
                            </div>
    
                           <!-- {{if eq $.Post.User $.User.Name}}
                            <div>
                                <a class="btn submitbtn" href="/edit/{{.ID}}">Edit</a>
                                <a class="btn cancelbtn" href="/delete/{{.ID}}" data-method="delete"
                                    data-confirm="Are you sure?">Delete</a>
    
                            </div>
                            {{end}}-->
    
                        
                        
                        </div>
    
                        <div class="postimgdiv">
                            {{with .ImageSrc}}
                            <img class="postimg" src="/static/{{.}}" alt="">
                            {{end}}
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
                            {{range .Comments}}
                            <div class="insidecatepadding">
                                <div class="catecard">
                                    <div class="card-body">
                                        <div class="post-card" id="comment{{.ID}}">
                                            <div class="votingbuttons">
                                                </form>
                                                <div class="post-upvote">
                                                    <form id="comment-up-{{.ID}}"
                                                        onsubmit="return fetchcomment('{{.ID}}', 'up')">
                                                        <input value="{{.ID}}" hidden>
                                                        <button type="submit" class="btn up-button"
                                                            onclick="vote('comment{{.ID}}', 'up-button')">
                                                            <i class="bi upb bi-caret-up"></i>
                                                        </button>
                                                    </form>
                                                </div>
    
                                                <div class="votecounter">
                                                    <a name="votes" class="votes" value="0">{{.Votes}}</a>
                                                </div>
    
                                                <div class="post-downvote">
                                                    <form id="comment-down-{{.ID}}"
                                                        onsubmit="return fetchcomment('{{.ID}}', 'down')">
                                                        <input value="{{.ID}}" hidden>
                                                        <button class="btn down-button"
                                                            onclick="vote('comment{{.ID}}', 'down-button')">
                                                            <i class="bi downb bi-caret-down"></i>
                                                        </button>
                                                    </form>
                                                </div>
                                               
                                            </div>
    
                                            <div class="postdetails">
                                                <div class="post-username">Posted by {{.User}} {{timeAgo .Created}}</div>
                                                <p class="post-description-comment card-text">{{.Content}}</p>
                                            </div>
                                         <!--   {{if eq .User $.User.Name}}
                                            <div>
                                                <div>
                                                    <a class="btn submitbtn" href="" data-toggle="modal" data-target="#exampleModal">Edit</a>
                                                    <a class="btn cancelbtn" href="/deletecomment/{{.ID}}">Delete</a>
                                                </div>
                                            </div>
                                            {{end}}-->
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
                        </div>							
                            {{end}}
                        </div>
                    </div>
    
                </div>
    
    
            </div>
        </div>
    </div>
        `;
    }
}