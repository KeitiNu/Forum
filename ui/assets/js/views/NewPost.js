import AbstractView from "./AbstractView.js";


export default class extends AbstractView {
    constructor(params) {
        super(params);
        this.setTitle("Kodify - Post");
    }


	


    async getHtml() {
        return `
	<div style="display:none">${this.doSubmit}</div>
    <div class="createpostcontent">
	<div class="createpostbox">
		<div class="createpostboxinside">
			<div class="createpostheader">
				<h1>Create a post</h1>
			</div>
			<form method="POST" id="threadform" enctype="multipart/form-data">
				<!-- Thread title -->
				<div class="">
					<label class="form-label labels">Title</label>
					<input class="form-control" type="text" name="title" placeholder="Post title" maxlength="40" />
					<label class='error' id='errortitle'></label>
				</div>

				<!-- Content textarea -->
				<div class="">
					<label class="form-label labels">Description</label>
					<textarea class="form-control" id="postcontent" name="content" onkeyup="charactercount()"
						maxlength="500" placeholder="Maximum 500 characters"></textarea>
					<div class="d-flex justify-content-end">
						<span id="words_count" class="badge badger">
							<span id="textcount">0</span>/500
						</span>
					</div>
				</div>
				<div id="textarea_count" class="badge pull-right"></div>
				<label class='error' id='errorcontent'></label>

				<!-- Category tagging -->
				<div class="">
					<label class="form-label labels">Tags</label>
					<fieldset>
					${this.params.Categories.map(function(cat){
						return `<label style="color:white;"><input type="checkbox" name="category" value=`+cat.Title+`>
						`+cat.Title+`</label>`
						}).join("")}
					</fieldset>
					<label class='error' id='errorcategory'></label>
				</div>

				<!-- Sumbit button -->
				<div class="submitpostbuttons">
					<a href="/" class="btn cancelbtn">Cancel</a>
					<button type="submit" class="btn submitbtn" name="submitPost">Submit button</button>
				</div>
			</form>
		</div>
	</div>
</div>

    `;
    }
}


