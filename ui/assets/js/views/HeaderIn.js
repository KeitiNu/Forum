import AbstractView from "./AbstractView.js";

export default class extends AbstractView {
    constructor(params) {
        super(params);
    }

    async getHtml() {
        return `
    <nav class="navbar header py-1 navbar-expand-sm fixed-top">
        <div class="container-fluid">
          <a class="navbar-brand header" href="/" data-link>
            <img id="header-mainlogo" src="static/css/images/logo_white.png" class="rounded-pill"></a>
                <form class="d-flex">
                    <div class="dropdown header">
                        <button class="btn inheader btn-secondary dropdown-toggle" type="button" id="dropdownMenuButton2" data-bs-toggle="dropdown" aria-expanded="false">
LOGGED IN
                        </button>
                        <ul class="dropdown-menu header dropdown-menu-dark" aria-labelledby="dropdownMenuButton2">
                          <li><a class="dropdown-item header bi-person-circle" href="/profile" data-link>   Profile</a></li>
                          <li><hr class="dropdown-divider"></li>
                          <li><a class="dropdown-item header bi-box-arrow-in-right" href="/logout" data-link>   Log Out</a></li>
                        </ul>
                      </div>
                </form>
        </div>
    </nav>
        `;
    }
}