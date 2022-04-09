// window.addEventListener('hashchange', () => {
//     console.log(location.hash)
// })


// const express = require('express')
// const path = require('path')
// const app = express()

// // connecting to the static folder so we can connect with index.js
// app.use('/static', express.static(path.resolve(__dirname, 'frontend', 'static')))

// app.get('/*', (req, res) => {
//     res.sendFile(path.resolve(__dirname, 'frontend', 'index.html'))
// })

// // starting server at 8090
// app.listen(process.env.POST || 8090, () => {
//     console.log('Server is running')
// })