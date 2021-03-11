const express = require('express')
const router = express.Router()

// Add your routes here - above the module.exports line
router.get('*', (req, res, next) => {
  req.session.data.url = req.url
  next()
})

router.get('/home', (req, res, next) => {
  if (req.session.data.email === 'user@opgtest.com') {
    res.redirect('/my-details')
  } else {
    next()
  }
})

module.exports = router
