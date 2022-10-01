// const sqlite3 = require('sqlite3').verbose();
let user = undefined
let recipient = undefined

/* Opening and loading old messages */

let disable = false;
let offset = 0
// Open chat between two two users
// also applying the event listeners 
const openChat = async (e) => {
    let activity = document.getElementById('activity')
    let dialog = document.getElementById('dialog')
    let chat = document.getElementById("chat_area")

    let input = {
        "container": document.getElementById('input'),
        "input": document.getElementById("input_text"),
        "button": document.getElementById("input_button")
    }

    applyEventListeners()

    // from the event we can get the id of the element
    // that was just clicked. Based on that we will open the chat
    console.log("id:\n", e)
    if (!disable) {
        disable = true

        // user = e.currentTarget
        recipient = e.currentTarget.getAttribute("data-username")
        user = e.currentTarget.getAttribute("data-currentuser")

        let recipientField = document.getElementById('recipientId')
        recipientField.value = recipient;

        document.getElementById("input_text").textContent = recipient
        offset = 0

        if (input.container.classList.length === 1) {
            await collapse(activity, dialog, input)
        }

        var resp = await fillChatLog(user, recipient, offset)



        console.log(resp)
        resp.forEach((message) => {
            let reciever = message.Recipient == recipient ? "recipient" : "user";
            let username = message.Recipient == recipient ? recipient : user;
            let bubble = createBubble(message.Content, username, reciever, message.SentAt);
            chat.insertBefore(bubble, chat.firstChild)
        })
        extend(activity, dialog, input)
    }
}

const applyEventListeners = async () => {
    // When the user presses enter, we send the message
    let textBox = document.getElementById('input_text');
    if (textBox) {
        textBox.addEventListener("keydown", (e) => {
            if (e.code === "Enter") send()
        })
    }

    // When the user scrolls up in the dialog box, we load 10 more messages
    let scrollListen = document.getElementById('chat_area')
    if (scrollListen) {
        scrollListen.addEventListener("scroll", async (elem) => {
            if (elem.target.scrollTop === 0) {
                let resp = await fillChatLog(user, recipient, offset);
                offset += 10

                resp.forEach((message) => {
                    let reciever = message.Recipient == recipient ? "recipient" : "user";
                    let username = message.Recipient == recipient ? recipient : user;
                    let bubble = createBubble(message.Content, username, reciever, message.SentAt);
                    elem.target.insertBefore(bubble, elem.target.firstChild)
                })
            }
        })
    }
}


async function fillChatLog(user, recipient, offset) {
    console.log(user, recipient, offset)

    var values =             {
        User: user,
        Recipient: recipient,
        Offset: offset,
    }

    var obj = await fetch('/message', {
        method: 'POST',
        headers: {
            'Content-type': 'application/json; charset=UTF-8'
        },
        body: JSON.stringify(values)
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

const collapse = async (activity, dialog, input) => {
    if (input.input.classList.length === 1) {
        addClass(input.container, "none")
        addClass(input.input, "none")
        addClass(input.button, "none")
    } else {
        changeClass(input.container, "none")
        changeClass(input.input, "none")
        changeClass(input.button, "none")
    }
    changeClass(activity, "extended")
    changeClass(dialog, "none")

    await delay(400)   
    changeClass(input.input, "remove")
    changeClass(input.button, "remove")
    
    await delay(600)
    changeClass(input.container, "remove")
    changeClass(dialog, "remove")
    document.getElementById("chat_area").innerHTML = ""
}

const extend = async (activity, dialog, input) => {
    changeClass(activity, 'minimized')
    changeClass(dialog, 'none')
    removeClass(input.container)

    await delay(10)
    changeClass(dialog, 'large')

    await delay(490)
    removeClass(input.button)
    removeClass(input.input)
    
    await delay(500)
    disable = false
}

function delay(milliseconds){
    return new Promise(resolve => {
        setTimeout(resolve, milliseconds);
    });
}

const changeClass = (elem, value) => {
    let arr = elem.className.split(' ')
    arr[arr.length-1] = value
    elem.className = arr.join(' ')
}
const addClass = (elem, value) => {
    let arr = elem.className.split(' ')
    arr.push(value)
    elem.className = arr.join(' ')
}
const removeClass = elem => {
    let arr = elem.className.split(' ')
    arr.pop()
    elem.className = arr.join(' ')
}

/* Sending messages to the chat */
const send = () => {
    let input = document.getElementById('input_text')
    let bubble = createBubble(input.value, user, "user", new Date())
    document.getElementById("chat_area").appendChild(bubble)
}

const createBubble = (text, name, style, time) => {
    let container = document.createElement('div')
    let info = document.createElement('div')
    let username = document.createElement('p')
    let date = document.createElement('p')
    let bubble = document.createElement('p')

    container.className = `${style}`
    info.className = "info"
    bubble.className = "bubble"
    username.textContent = name
    date.textContent = time
    bubble.textContent = text

    info.appendChild(username)
    info.appendChild(date)
    container.appendChild(info)
    container.appendChild(bubble)

    return container
}



// status: what status do you want th user to be changed to 
// 0: offline   1:oneline
const changeStatus = (username, status) => {
    let div = document.getElementById(`status-${username}`)

    if (div!= null) {
        if (status == 0 && div.classList.length == 1) {
            addClass(div, "away")
        } else if (status == 1 && div.classList.length == 2) {
            removeClass(div)
        }
    }
}