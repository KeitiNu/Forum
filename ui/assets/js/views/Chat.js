import AbstractView from "./AbstractView.js";

export default class extends AbstractView {
    constructor(params) {
        super(params);
    }
    

  

    async getHtml() {
        let users = this.params.Users
        if (!users) {users = []}

        return`
        <div style="display:none">${this.sendMessage}</div>
        <div class="chat">
            <div id="activity" class="scroll box extended">
                <div id="activity" class="activity">
                    <div id="Keiti" class="user away" onclick="openChat(event)">
                        <span class="status"></span>
                        <p class="name">Keiti</p>
                    </div>
                    
                    ${
                        users.map(function(user) {
                        return `
                        <div id="status-${user.Name}" class="user away" onclick="openChat(event)">
                            <span class="status"></span>
                            <p class="name">${user.Name}</p>
                        </div>
                        `
                    }).join("")
                }
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

                <div id="input" class="input remove" method="POST">
                    <input id="input_text" class="quick remove" name="message" type="text" name="Message"></input>
                    <input id="RecipientId" style="display: none"></input>
                    <input style="display: none" name="UserId">${this.params.AuthenticatedUser.Name}</input>
                    <button id="input_button" data-userId="" class="quick remove" onClick="send()">Send!</button>
                </div>
            </div>
        </div>
        `
    }
}