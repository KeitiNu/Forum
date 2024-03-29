import AbstractView from "./AbstractView.js";

export default class extends AbstractView {
    constructor(params) {
        super(params);
        this.setTitle("Kodify");
    }
    
    async getHtml() {
        return `<div class="mainpagecontent">
        <div class="mainpagebox">
            <div class="mainpageboxinside">
                    <div class="mainpageboxheader">
                        <h1>Hello, <a class="kodify" href="/profile" data-link>${this.params.AuthenticatedUser?this.params.AuthenticatedUser.Name:""}</a>! Welcome to <a class="kodify" href="/" data-link>kodify</a>!</h1>
                        <div class="adderbuttons">
                            <a class="btn adderbutton" href="/submit" data-link>Create Post</a>
                        </div>
                    </div>
                    <div class="categories">
                        <div class="insidecategories">
            
                            ${this.params.Categories.map(function(par){

                            return `<div>
                                        <div class="insidecatepadding">
                                            <div class="catecard">
                                                <div class="card-body">
                                                    <a href="/category/`+ par.Title+ `" data-link class="catecardtitle stretched-link">`+ par.Title+ `</a> 
                                                    <p class="card-text">`+par.Description+`</p>
                                                </div>
                                            </div>
                                        </div>
                                    </div>`
                            })}

                        </div>
                    </div>
                </div>
            </div>
        </div>`;
    }

    
}
