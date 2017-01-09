const express = require('express')
const bcrypt = require('bcrypt')

const BCRYPT_ROUNDS = 10

module.exports = function(database) {
  const router = express.Router()

  //
  // Authorisation
  //

  function authorise(req, res, next) {
    if (req.session.card) {
      return next()
    }

    let err = new Error('Must log in')
    err.status = 401

    return next(err)
  }


  //
  // Index
  //

  router.get('/', function(req, res, next) {
    res.render('index', {})
  })


  //
  // Register
  //

  router.get('/register', function(req, res, next) {
    res.render('register', {
      didRegister: false
    })
  })

  router.post('/register', function(req, res, next) {
    const card = (req.body.card_id || '').trim()
    const pass = (req.body.card_pass || '')

    if (card == '' || pass == '') {
      let err = new Error('Invalid Card ID or password')
      err.status = 400

      return next(err)
    }

    bcrypt.hash(pass, BCRYPT_ROUNDS, function(err, hash) {
      if (err != null) {
        return next(err)
      }

      database.query('INSERT INTO cards VALUES ($1, $2)', [card, hash], function(err) {
        if (err != null) {
          return next(err)
        }

        res.render('register', {
          didRegister: true
        })
      })
    })
  })


  //
  // Login
  //

  router.get('/login', function(req, res, next) {
    res.render('login', {
      attemptedLogin: false
    })
  })

  router.post('/login', function(req, res, next) {
    const card = (req.body.card_id || '').trim()
    const pass = (req.body.card_pass || '')

    if (card == '' || pass == '') {
      let err = new Error('Invalid Card ID or password')
      err.status = 400

      return next(err)
    }

    database.query('SELECT id, password FROM cards WHERE id = $1', [card], function(err, result) {
      if (err != null) {
        return next(err)
      }

      if (result.rowCount != 1) {
        let err = new Error('Invalid login')
        err.status = 401

        return next(err)
      }

      const hash = result.rows[0].password
      bcrypt.compare(pass, hash, function(err, success) {
        if (err != null) {
          return next(err)
        }

        if (success === false) {
          return res.render('login', {
            attemptedLogin: true
          })
        }

        req.session.card = card
        res.redirect('/list')
      })
    })
  })


  //
  // List
  //

  router.get('/list', authorise, function(req, res, next) {
    const cardId = req.session.card

    database.query('SELECT id, list FROM cards WHERE id = $1', [cardId], function(err, result) {
      if (err != null) {
        return next(err)
      }

      if (result.rowCount != 1) {
        // something messed up in the database. force logout
        req.session.destroy()

        let err = new Error('Invalid number of rows returned')
        err.status = 500

        return next(err)
      }

      const list = result.rows[0].list
      res.render('list', {
        list: list
      })
    })
  })


  //
  // List management
  //

  router.get('/list/add', authorise, function(req, res, next) {
    res.render('add', {})
  })

  router.post('/list/add', authorise, function(req, res, next) {
    const cardId = req.session.card
    const name = (req.body.item_name || '').trim()

    database.query('UPDATE cards SET list = array_append(list, $1) WHERE id = $2', [name, cardId], function(err) {
      if (err != null) {
        return next(err)
      }

      res.redirect('/list')
    })
  })

  return router
}
