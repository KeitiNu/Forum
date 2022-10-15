let user = undefined
let recipient = undefined

let disable = false;
let offset = 0

let removeTimer = 0;

/* Opening and loading old messages */

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
    if (!disable) {
        
        let recipientField = document.getElementById('recipientId')
        recipient = e.currentTarget.getAttribute("data-username")
        user = e.currentTarget.getAttribute("data-currentuser")
        
        let bell = document.getElementById(`bell-${recipient}`)
        recipientField.value = recipient;
        // let input_field = document.getElementById("input_text")
        document.getElementById("input_text").textContent = recipient
        
        $('#messageDiv').on('keyup', '#input_text', function () {
            debounce(typingInProgress(user, recipient), 2000, true)
        });

        disable = true
        offset = 0

        if (input.container.classList.length === 1) {
            await collapse(activity, dialog, input)
        }
        if (bell.classList.length == 2) {
            removeClass(bell)
        }

        var resp = await fillChatLog(user, recipient, offset)
        offset += 10;

        resp.forEach((message) => {
            let bubble = addMessage(message)
            chat.insertBefore(bubble, chat.firstChild)
        })

        extend(activity, dialog, input)
    }
}

const applyEventListeners = async () => {
    // When the user scrolls up in the dialog box, we load 10 more messages
    let chat = document.getElementById('chat_area')
    if (chat) {
        chat.addEventListener("scroll", async (elem) => {
            if (elem.target.scrollTop === 0) {
                let height = chat.scrollHeight
                let resp = await fillChatLog(user, recipient, offset);
                offset += 10

                resp.forEach((message) => {
                    let bubble = addMessage(message)
                    elem.target.insertBefore(bubble, elem.target.firstChild)
                    chat.scrollTop = height
                })
            }
        })
    }
}


async function fillChatLog(user, recipient, offset) {
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
    let chat = document.getElementById("chat_area")
    
    changeClass(activity, 'minimized')
    changeClass(dialog, 'none')
    removeClass(input.container)
    
    chat.scrollTop = chat.scrollHeight
    await delay(10)
    changeClass(dialog, 'large')
    
    chat.scrollTop = chat.scrollHeight
    await delay(490)
    removeClass(input.button)
    removeClass(input.input)
    
    chat.scrollTop = chat.scrollHeight
    await delay(500)
    disable = false
    chat.scrollTop = chat.scrollHeight
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

// Takes in a date and return the correct message format
const getDateformat = (date) => {
    let today = new Date()

    if (today.getFullYear() == date.getFullYear() &&
        today.getMonth() == date.getMonth() &&
        today.getDate() == date.getDate()) {
        let hours = date.getHours() < 10 ? `0${date.getHours()}` : date.getHours()
        let minutes = date.getMinutes() < 10 ? `0${date.getMinutes()}` : date.getMinutes()
        return `${hours}:${minutes}`
    } else {
        return `${date.getDate()}/${date.getMonth()}/${date.getFullYear()}`
    } 
}

/* Sending messages to the chat */
const addMessage = (message) => {
    let date = getDateformat(new Date(message.SentAt))
    let receiver = message.Recipient == recipient ? "user" : "recipient";
    let username= message.Recipient == recipient ? user : recipient;
    return createBubble(message.Content, username, receiver, date);
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
// 0: offline   1:online
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

const notify = (sender, message) => {
    if (sender == recipient) {
        let chat = document.getElementById("chat_area")
        let bubble = createBubble(message, sender, "recipient", getDateformat(new Date()))
        chat.appendChild(bubble)
        chat.scrollTop = chat.scrollHeight
    } else {
        let div = document.getAnimations(`status-${sender}`)
        let bell = document.getElementById(`bell-${sender}`)
        if (bell.classList.length == 1) {
            addClass(bell, "notif")
        }
        moveToTop(sender)
    }
}

const typing = (sender) => {
    let chat = document.getElementById("chat_area")
    var loadingDiv = chat.lastElementChild;
    var alreadyloading = loadingDiv != null?  loadingDiv.classList.contains("loading"):false
    // var alreadyLoading = chat.lastElementChild.classList.contains("loading");

    if (sender == recipient && !alreadyLoading && sender != user) {

        let container = document.createElement('div')
        let info = document.createElement('div')
        let username = document.createElement('p')

        container.className = "recipient loading"
        info.className = "info"
        username.textContent = sender

        let img = document.createElement("img");
        img.src = 'static/css/images/dots.gif'
        img.width = 24

        info.appendChild(username)
        info.appendChild(img)
        container.appendChild(info)
        chat.appendChild(container)
        startRemoveCouter();
    }else{
        resetRemoveCouter()
    }
}

function resetRemoveCouter(){
    clearTimeout(removeTimer);
    startRemoveCouter();

}

function startRemoveCouter(){
    console.log("here")
     removeTimer = setTimeout(() => {
        $( ".loading" ).remove();
    }, 1000);
}

// moves the user to the top of activities list
const moveToTop = (username) => {
    let activity = document.getElementById("activity")
    let divs = activity.childNodes
    let index = 0

    divs.forEach((elem, i) => {
        if (elem.id == `status-${username}`) {
            index = i
        }
    })

    divs.unshift(data.splice(index, 1)[0]);
    activity.childNodes = divs
}



function typingInProgress(sender, recipient){

    var values =             {
        UserId: sender,
        RecipientId: recipient,
    }

    var obj = fetch('/typing', {
        method: 'POST',
        headers: {
            'Content-type': 'application/json; charset=UTF-8'
        },
        body: JSON.stringify(values)
    })
        .then(response => {
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



function debounce(func, wait, immediate) {
    var timeout;
    return function() {
        var context = this,
            args = arguments;
        var later = function() {
            timeout = null;
            if (!immediate) func.apply(context, args);
        };
        var callNow = immediate && !timeout;
        clearTimeout(timeout);
        timeout = setTimeout(later, wait);
        if (callNow) func.apply(context, args);
    };
};