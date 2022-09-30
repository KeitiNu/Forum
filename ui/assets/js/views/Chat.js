import AbstractView from "./AbstractView.js";

export default class extends AbstractView {
    constructor(params) {
        super(params);
    }

    async getHtml() {
        return`

        <div class="chat">
            <div id="activity" class="scroll box extended">

                <div id="activity" class="activity">

                    <div id="Laura-Eliise" class="user" onclick="openChat(event)">
                        <span class="status active"></span>
                        <p class="name">Laura-Eliise</p>
                    </div>

                    <div id="Keiti" class="user away" onclick="openChat(event)">
                        <span class="status"></span>
                        <p class="name">Keiti</p>
                    </div>
                    
                </div>

            </div>

            <div id="dialog" class="dialog box remove">
                <div class="container" id="chat_area">
                    <div class="user">
                        <div class="info">
                            <p>Laura-Eliise</p>
                            <p>21:30</p>
                        </div>
                        <p class="bubble">Hello!</p>
                    </div>
                    <div class="recipient">
                        <div class="info">
                            <p>Keiti</p>
                            <p>21:33</p>
                        </div>
                        <p class="bubble">Hey! Working hard or hardly working!</p>
                    </div>
                </div>

                <div id="input" class="input remove">
                    <input id="input_text" class="quick remove" type="text"></input>
                    <button id="input_button" class="quick remove" onClick="send()">Send!</button>
                </div>
            </div>
        </div>
        `
    }
}