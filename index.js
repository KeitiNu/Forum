// const navigateTo = url => {
//     history.pushState(null, null, url)
//     router()
// }

// const router = async () => {
//     // list of all the pages the user can visit
//     const routes = [
//         {path: '/', view: () => console.log('Viewing Dashboard')},
//         {path: '/posts', view: () => console.log('Viewing Posts')},
//         {path: '/settings', view: () => console.log('Viewing Settings')},
//     ]

//     // Test each route for potential match
//     const potentialMatches = routes.map(route => {
//         return {
//             route: route,
//             isMatch: location.pathname === route.path
//         }
//     })

//     let match = potentialMatches.find(potentialMatches => potentialMatches.isMatch)

//     if (!match) {
//         return {
//             route: routes[0],
//             isMatch: true
//         }
//     }
// }

// // makes it so that the user can go back to a previously visited page
// window.addEventListener('popstate', router)

// // starts the router as soon as the page is loaded and
// // adds an event listener for when a link is pressed that redirects
// // the user to a new page
// document.addEventListener('DOMContentLoaded', () => {
//     document.body.addEventListener('click', e => {
//         if (e.target.matches('[data-link]')) {
//             e.preventDefault()
//             navigateTo(e.target.href)
//         }
//     })
//     router()
// })