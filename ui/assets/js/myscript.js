// function vote(id, updown, postcom) {
//     let post = document.getElementById(id);
//     let button = post.getElementsByClassName(updown)[0];
//     let vote = post.getElementsByClassName("votes")[0];
//     let curVotes = parseInt(vote.innerHTML);
//     let number = 1
//     let otherbutton = post.getElementsByClassName('down-button')[0];

//     if (updown == 'down-button') {
//         number = -1;
//         otherbutton = post.getElementsByClassName('up-button')[0];
//     };

//     if (button.classList.contains("active")) {
//         number *= -1;
//         button.classList.remove("active");
//     } else {
//         button.classList.add("active");
//     }

//     if (otherbutton.classList.contains("active")) {
//         otherbutton.classList.remove("active");
//         vote.innerHTML = curVotes + (number * 2);
//     } else {
//         vote.innerHTML = curVotes + number
//     }
// }

// function fetchpost(id, updown) {
//     // (A) GET FORM DATA
//     // var form = document.getElementById("form-"+updown+"-"+id);
//     var formData = new FormData();
//     formData.append("postID", id);
//     formData.append("type", updown);

//     // (B) FETCH
//     fetch("http://localhost:8090/test", {
//         method: "POST",
//         body: formData,
//     })
//         .then((res) => { return res.text(); })
//         .then((txt) => { console.log(txt); })
//         .catch((err) => { console.log(err); });

//     // (C) PREVENT HTML FORM SUBMIT
//     return false;
// }

// function fetchcomment(id, updown) {
//     // (A) GET FORM DATA
//     // var form = document.getElementById("form-"+updown+"-"+id);
//     var formData = new FormData();
//     formData.append("postID", id);
//     formData.append("type", updown);

//     // (B) FETCH
//     fetch("http://localhost:8090/testcomment", {
//         method: "POST",
//         body: formData,
//     })
//         .then((res) => { return res.text(); })
//         .then((txt) => { console.log(txt); })
//         .catch((err) => { console.log(err); });

//     // (C) PREVENT HTML FORM SUBMIT
//     return false;
// }


function charactercount() {
    // const areatextarea = document.querySelector("#postcontent");
    const areatext = document.querySelector("#postcontent").value.length;
    const textcount = document.querySelector("#textcount");
    // const wordcount = document.querySelector("#words_count");
    textcount.innerHTML = areatext;
};


// function myFunction() {
//     var x = document.getElementById("myDIV");
//     if (x.style.display === "none") {
//         x.style.display = "block";
//     } else {
//         x.style.display = "none";
//     }
// };



// function signInButtonPressed(e) {
//     e.preventDefault();
//     const data = new FormData(e.target);
//     const value = Object.fromEntries(data.entries());

// var o = await abc(value)

// console.log(o)

// }



// async function abc(value){

//     var obj =  fetch('/data/login', {
//         method: 'POST',
//         headers: {
//             'Content-type': 'application/json; charset=UTF-8'
//         },
//         body: JSON.stringify(value)
//     })
//     .then(response => {
    
//         if (!response.ok) {
//             throw new Error(`HTTP error: ${response.status}`);
//         }
//         // Otherwise (if the response succeeded), our handler fetches the response
//         // as text by calling response.text(), and immediately returns the promise
//         // returned by `response.text()`.
//         return response.text()
    
//     })
//     .then(json => JSON.parse(json))
//     .catch(err => console.error(`Fetch problem: ${err.message}`))
    
    
    
//      return obj
//     }
