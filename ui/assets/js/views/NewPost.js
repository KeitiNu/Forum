import AbstractView from "./AbstractView.js";


export default class extends AbstractView {
    constructor(params) {
        super(params);
        this.setTitle("Kodify - Post");
    }



	get doSubmit(){

		$(document.body).on('submit', 'form', async function (e) {
        e.preventDefault();

        var data = new FormData(e.target);
        const cats = data.getAll('category');

        var values = Object.fromEntries(data.entries());
        values.category = cats

        const location = window.location.pathname
        var o = await fetchFormData(values, location)

		this.params = o

			const errors = this.params.Form.Errors.Errors
			const keys = Object.keys(errors)

			if (keys.length == 0) {
				const tempLink = document.createElement('a')
				const tempLocation = document.querySelector('.submitpostbuttons')

				tempLink.href = '/post/'+this.params.Sort
				tempLink.dataset.link

				tempLocation.appendChild(tempLink)
				tempLink.click()

			}else{
				var errorSpots = document.querySelectorAll('.error')

				errorSpots.forEach(err => {
					err.innerHTML = ""
				});

				keys.map(function(key){
					var spot = $('#error'+key)
					spot.text(errors[key])
				})
			}

    });


	async function fetchFormData(value, url) {

		var obj = fetch('/data'+url, {
			method: 'POST',
			headers: {
				'Content-type': 'application/json; charset=UTF-8'
			},
			body: JSON.stringify(value)
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

	}




    async getHtml() {
        return `
		${this.doSubmit}
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
					<label class="form-label labels">Tags (maximum 5)</label>
					<fieldset>
					${this.params.Categories.map(function(cat){
						return `<input type="checkbox" name="category" value=`+cat.Title+`>`+cat.Title+`</input>`
						})}
				<!--	<select class="form-select form-control" id="tags-input" name="category" multiple
						data-allow-clear="true" data-show-all-suggestions="true" data-suggestions-threshold="0"
						data-max="5">


						{{range .Categories}}
						<option class="dropdownoption" value="{{.Title}}">{{.Title}}</option>
						{{end}}
					</select>-->
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


