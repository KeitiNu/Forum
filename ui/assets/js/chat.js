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
        recipient = e.target.id != "" ? e.target.id : recipient

        document.getElementById("input_text").textContent = recipient

        if (input.container.classList.length === 1) {
            await collapse(activity, dialog, input)
            await fillLog(recipient, 0)
            extend(activity, dialog, input) 
        } else {
            extend(activity, dialog, input)
        }
    }
}

const applyEventListeners = () => {
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
        scrollListen.addEventListener("scroll", (elem ,_) => {
            if (elem.target.scrollTop === 0) {
                /* Query 10 messages from the database */
                let bubble = createBubble('Test', "bot", 'recipient', "Too Late")
                elem.target.insertBefore(bubble, elem.target.firstChild)
            }
        })
    }
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
    if (input.value.length != 0) {
        let bubble = createBubble(input.value, "Laura-Eliise", "user", "23:00")
        input.value = ""
        document.getElementById("chat_area").appendChild(bubble)
    }
}

const fillLog = async (recipient, offset) => {
    let arr = ["Laura", "Laura-Eliise", "Keiti"]

    arr.forEach( (user) => {
        createBubble("test", user, recipient, "23:00")
    })
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



/* FUNCTIONS WE MAY NEED IN THE FUTURE */

const changeStatus = (username) => {
    let div = document.getElementById(`status-${username}`)
    
    if (div.classList.length == 1) {
        addClass(div, "away")
    } else {
        removeClass(div)
    }
}