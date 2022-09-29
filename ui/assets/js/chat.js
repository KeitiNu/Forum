let disable = false;

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
    console.log("id", e.target.id)
    if (!disable) {
        disable = true
        if (input.container.classList.length === 1) {
            collapse(activity, dialog, input)
            setTimeout(() => {
                disable = true
                extend(activity, dialog, input) 
            }, 1001)
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
    let scrollListen = document.getElementById('chatArea')
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

    setTimeout(() => {
        changeClass(input.container, "remove")
        changeClass(dialog, "remove")
    }, 1000)
    setTimeout(() => {
        changeClass(input.input, "remove")
        changeClass(input.button, "remove")
    }, 700)

    setTimeout(() => {disable = false}, 1001)
}

const extend = (activity, dialog, input) => {
    changeClass(activity, 'minimized')
    changeClass(dialog, 'none')
    removeClass(input.container)

    setTimeout(() => {
        changeClass(dialog, 'large')
    }, 10)
    setTimeout(() => {
        removeClass(input.button)
        removeClass(input.input)
    }, 500)

    setTimeout(() => {
        disable = false
    }, 1001)
}

const changeClass = (elem, value) => {changeClass
    let arr = elem.className.split(' ')
    arr[arr.length-1] = value
    elem.className = arr.join(' ')
}
const addClass = (elem, value) => {
    let arr = elem.className.split(' ')
    arr.push(value)
    elem.className = arr.join(' ')
}
const removeClass = elem => {changeClass
    let arr = elem.className.split(' ')
    arr.pop()
    elem.className = arr.join(' ')
}


const send = () => {
    let input = document.getElementById('input_text')
    if (input.value.length != 0) {
        let bubble = createBubble(input.value, "Laura-Eliise", "user", "23:00")
        input.value = ""
        document.getElementById("chatArea").appendChild(bubble)
    }
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