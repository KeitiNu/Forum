let votes = document.querySelector('#uservotes').innerText;
votes = votes.split(`]`).slice(0, 4);
const upvotes = votes[0].slice(2).split(' ');
const downvotes = votes[1].slice(2).split(' ');
const commentUpvotes = votes[2].slice(2).split(' ');
const commentDownvotes = votes[3].slice(2).split(' ');
console.log('upvotes: ', upvotes, 'downvotes: ', downvotes,  'commentupvotes: ', commentUpvotes, 'commentdownvotes: ', commentDownvotes);
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
for (const comment of commentUpvotes) {
    let button = document.querySelector(`#comment-up-${comment} button`);
    if (button) {
        button.classList.add("active");
    }
}
for (const comment of commentDownvotes) {
    let button = document.querySelector(`#comment-down-${comment} button`);
    if (button) {
        button.classList.add("active");
    }
}
function vote(id, updown, postcom) {
    let post = document.getElementById(id);
    let button = post.getElementsByClassName(updown)[0];
    let vote = post.getElementsByClassName("votes")[0];
    let curVotes = parseInt(vote.innerHTML);
    let number = 1
    let otherbutton = post.getElementsByClassName('down-button')[0];

    if (updown == 'down-button') {
        number = -1;
        otherbutton = post.getElementsByClassName('up-button')[0];
    };

    if (button.classList.contains("active")) {
        number *= -1;
        button.classList.remove("active");
    } else {
        button.classList.add("active");
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

function fetchcomment(id, updown) {
    // (A) GET FORM DATA
    // var form = document.getElementById("form-"+updown+"-"+id);
    var formData = new FormData();
    formData.append("postID", id);
    formData.append("type", updown);

    // (B) FETCH
    fetch("http://localhost:8090/testcomment", {
        method: "POST",
        body: formData,
    })
        .then((res) => { return res.text(); })
        .then((txt) => { console.log(txt); })
        .catch((err) => { console.log(err); });

    // (C) PREVENT HTML FORM SUBMIT
    return false;
}