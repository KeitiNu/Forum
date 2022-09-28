let disable = false
let activity = document.getElementById('activity')
let dialog = document.getElementById('dialog')
let input = {
    "container": document.getElementById('input'),
    "input": document.getElementById("input_text"),
    "button": document.getElementById("input_button")
}

// Open chat between these two users
const openChat = (e) => {
    // from the event we can get the id of the element
    // that was just clicked. Based on that we will open the chat
    // e.target.id
    if (!disable) {
        if (document.getElementById('input').classList.length === 1) {
            collapse(activity, dialog, input)
        } else {
            extend(activity, dialog, input)
        }
    }
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

const collapse = (activity, dialog, input) => {
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

    disable = true
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

    disable = true
    setTimeout(() => {disable = false}, 1001)
}