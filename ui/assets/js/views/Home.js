import AbstractView from "./AbstractView.js";

export default class extends AbstractView {
    constructor(params) {
        super(params);
        this.setTitle("Kodify");
    }
    

    async logData(){

        for (let index = 0; index < this.params.length; index++) {
            console.log("This params: ", this.params[index])
        }
            
    }


    async getHtml() {

        var strStart = `<div class="mainpagecontent">
        <div class="mainpagebox">
            <div class="mainpageboxinside">
                    <div class="mainpageboxheader">
                        <h1>Hello, <a class="kodify" href="/profile" data-link>{{.AuthenticatedUser.Name}}</a>! Welcome to <a class="kodify" href="/" data-link>kodify</a>!</h1>
                        <div class="adderbuttons">
                            <a class="btn adderbutton" href="/submit" data-link>Create Post</a>
                        </div>
                    </div>
                    <div class="categories">
                        <div class="insidecategories">`;

        var strMiddle = ``;

        for (let index = 0; index < this.params.length; index++) {
            strMiddle+= `
            
            <div>
            <div class="insidecatepadding">
                <div class="catecard">
                    <div class="card-body">
                               <a href="/category/${this.params[index].Title}" data-link class="catecardtitle stretched-link">${this.params[index].Title}</a> 
                        
                                <p class="card-text">${this.params[index].Description}</p>
                    </div>
                </div>
            </div>
        </div>
            
            `
        }

        var strEnd = `   </div>
        </div>
        </div>
        </div>
        </div>`

        return strStart+strMiddle+strEnd
        ;
    }


async getHtml2(){

var str = `<div class="mainpagecontent">
<div class="mainpagebox">
    <div class="mainpageboxinside">
            <div class="mainpageboxheader">
                <h1>Hello, <a class="kodify" href="/profile" data-link>{{.AuthenticatedUser.Name}}</a>! Welcome to <a class="kodify" href="/" data-link>kodify</a>!</h1>
                <div class="adderbuttons">
                    <a class="btn adderbutton" href="/submit" data-link>Create Post</a>
                </div>
               
            </div>
            <div class="categories">
                <div class="insidecategories">

                
                    {{range .Categories}}
                <div>
                    <div class="insidecatepadding">
                        <div class="catecard">
                            <div class="card-body">
                                       <a href="/category/${this.data.ID}" data-link class="catecardtitle stretched-link">${this.data.Title}</a> 
                                        <p class="card-text">${this.data.Description}</p>
                            </div>
                        </div>
                    </div>
                </div>
                {{end}}
                </div>
            </div>
    </div>
</div>
</div>`

return str

}








    
}
