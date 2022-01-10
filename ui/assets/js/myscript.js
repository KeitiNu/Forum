let votes = document.querySelector('#uservotes').innerText;
votes = votes.split(`]`).slice(0, 2);
const upvotes = votes[0].slice(2).split(' ');
const downvotes = votes[1].slice(2).split(' ');
for (const post of upvotes) {
    let button = document.querySelector(`#form-up-${post} button`);
    if (button) {
        button.classList.add("active");
    }
}
for (const post of downvotes) {
    let button = document.querySelector(`#form-down-${post} button`);
    if (button) {
        button.classList.add("active");
    }
}
function vote(id, which) {
    let post = document.getElementById(id);
    let button = post.getElementsByClassName(which)[0];
    let vote = post.getElementsByClassName("votes")[0];
    let curVotes = parseInt(vote.innerHTML);
    let number = 1
    let otherbutton = post.getElementsByClassName('down-button')[0];

    if (which == 'down-button') {
        number = -1;
        otherbutton = post.getElementsByClassName('up-button')[0];
    };

    if (button.classList.contains("active")) {
        number *= -1;
        button.classList.remove("active");
        console.log(button)
    } else {
        button.classList.add("active");
        console.log(button)
    }

    if (otherbutton.classList.contains("active")) {
        otherbutton.classList.remove("active");
        vote.innerHTML = curVotes + (number * 2);
    } else {
        vote.innerHTML = curVotes + number
    }
}

function fetchpost(id, updown) {
    // (A) GET FORM DATA
    // var form = document.getElementById("form-"+updown+"-"+id);
    var formData = new FormData();
    formData.append("postID", id);
    formData.append("type", updown);

    // (B) FETCH
    fetch("http://localhost:8090/test", {
        method: "POST",
        body: formData,
    })
        .then((res) => { return res.text(); })
        .then((txt) => { console.log(txt); })
        .catch((err) => { console.log(err); });

    // (C) PREVENT HTML FORM SUBMIT
    return false;
}