function charactercount() {
    // const areatextarea = document.querySelector("#postcontent");
    const areatext = document.querySelector("#postcontent").value.length;
    const textcount = document.querySelector("#textcount");
    // const wordcount = document.querySelector("#words_count");
    textcount.innerHTML = areatext;
};