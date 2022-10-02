import AbstractView from "./AbstractView.js";

export default class extends AbstractView {
    constructor(params) {
        super(params);
        this.setTitle("KodyFy - You weren't supposed to see this");
    }

    async getHtml() {
        return `
        <div class="errorcontentoutside">
            <div class="errorcontent">
                <div class="shade">
                    <div class="fourofour">
                    <h1>404</h1>
                    </div>
                    <div class="errortext">
                    <h2>Ooops! You weren't supposed to see this</h2>
                    </div>
                    <div class="returntext">
                    <h3>Go <a class="lowerh3a" href="/" data-link>BACK</a> and remember: you haven't seen 3anything.</h3>
                    </div>
                </div>
            </div>
        </div>
        `;
    }
}




