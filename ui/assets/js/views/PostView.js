import AbstractView from "./AbstractView.js";

export default class extends AbstractView {
    constructor(params) {
        super(params);
        this.setTitle("Kodify - Post");
    }

    get submitComment(){

        $(document.body).on('submit', 'form', async function (e) {
            e.preventDefault();
    
            var data = new FormData(e.target);
            var values = Object.fromEntries(data.entries());
    
            const location = window.location.pathname
            var o = await fetchFormData(values, "/comment")
    
            this.params = o
            console.log(o)

                const errors = this.params.Form.Errors.Errors
                const keys = Object.keys(errors)
    
                if (keys.length == 0) {
                //     const tempLink = document.createElement('a')
                //     const tempLocation = document.querySelector('.registerlink')

                //     if(o.AuthenticatedUser != null){
                //      document.cookie = "auth=true;"
                //     }

                //     tempLink.href = '/'
                //     tempLink.dataset.link
    
                //     tempLocation.appendChild(tempLink)
                //     tempLink.click()
    
                }else{
                //     var errorSpots = document.querySelectorAll('.error')
    
                //     errorSpots.forEach(err => {
                //         err.innerHTML = ""
                //     });
    
                //     keys.map(function(key){
                //         var spot = $('#error'+key)
                //         spot.text(errors[key])
                //     })
                }
    
        });
    
    
        async function fetchFormData(value, url) {
    
            var obj = fetch('/data'+url, {
                method: 'POST',
                headers: {
                    'Content-type': 'application/json; charset=UTF-8'
                },
                body: JSON.stringify(value)
            })
                .then(response => {
                    console.log("RESPONSE:", response)

                    if (!response.ok) {
                        throw new Error(`HTTP error: ${response.status}`);
                    }
                    // Otherwise (if the response succeeded), our handler fetches the response
                    // as text by calling response.text(), and immediately returns the promise
                    // returned by `response.text()`.
                    return response.text()
    
                })
                .then(json => JSON.parse(json))
                .catch(err => console.error(`Fetch problem: ${err.message}`))
    
    
    
            return obj
        }
    
    };


    async getHtml() {
        return `
        ${this.submitComment}
        <div class="mainpagecontent">
        <div class="mainpagebox">
            <div class="mainpageboxinside">
                <div class="insidecateboxheader">
                    <div class="mainbpageboxheaderthread">
                        <div class="post-card" id="post${this.params.Post.ID}">
    
                            <div class="postdetails">
                                <div class="post-username-thread">Posted by ${this.params.Post.User} ${moment(this.params.Post.Created).format("DD.MM.YYYY HH:mm")}</div>
                                <div class="post-title-thread"><a class="post-title-thread stretched-link"
                                        href="/post/${this.params.Post.ID}" data-link>${this.params.Post.Title}</a></div>
                                <div class="post-description-thread">${this.params.Post.Content}</div>
                            </div>
                        </div>
                    </div>
                    <div class="categories">
                        <div class="commentingbox">
                            <form class="comment" method="POST">
                                <input type="hidden" id="postId" name="postId" value="${this.params.Post.ID}">
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