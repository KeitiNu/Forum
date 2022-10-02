import AbstractView from "./AbstractView.js";

export default class extends AbstractView {
    constructor(params) {
        super(params);
    }
    
    async getHtml() {
        let users = this.params.Users
        let currentUser = this.params.AuthenticatedUser.Name
        if (!users) {users = []}

        return`
        <div style="display:none">${this.sendMessage}</div>
        <div class="chat">
            <div id="activity" class="scroll box extended">
                <div id="inner_activity" class="activity">                
                    ${
                        users.map(function(user) {
                        return `
                        <div id="status-${user.Name}"  data-username="${user.Name}" data-currentuser="${currentUser}" class="user away" onclick="openChat(event)">
                                <span class="status"></span>
                                <iconify-icon class="bell" icon="bi:bell-fill" id="bell-${user.Name}"></iconify-icon>
                            <p class="name">${user.Name}</p>
                        </div>
                        `
                    }).join("")
                }
                </div>

            </div>

            <div id="dialog" class="dialog box remove">
                <div class="container" id="chat_area">
                </div>

                <div id="input" class="input remove" method="POST">
                <form>
                    <input id="input_text" class="quick remove" name="Message" value="" autocomplete="off"></input>
                    <input type="hidden" name="RecipientId" id="recipientId"></input>
                    <input type="hidden" name="UserId" value="${currentUser}"></input>
                    <button id="input_button" data-userId="" class="quick remove" type="submit">Send!</button>

                </form>
                </div>
            </div>
        </div>
        `
    }
}