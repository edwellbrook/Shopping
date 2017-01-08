const express = require('express')
const path = require('path')
const http = require('http')
const logger = require('morgan')
const cookieParser = require('cookie-parser')
const bodyParser = require('body-parser')
const session = require('express-session')
const postgres = require('pg-pool')

const index = require('./routes/index')
const customerRoutes = require('./routes/customer')
const staffRoutes = require('./routes/staff')

const app = express()
const database = postgres()

app.set('views', path.join(__dirname, 'views'))
app.set('view engine', 'ejs')

app.use(logger('dev'))
app.use(bodyParser.json())
app.use(bodyParser.urlencoded({ extended: false }))
app.use(cookieParser())
app.use(express.static(path.join(__dirname, 'public')))

app.use(session({
  secret: 'some-secret',
  saveUninitialized: false,
  resave: true
}))

app.use('/', index)
app.use('/customer', customerRoutes(database))
app.use('/staff', staffRoutes(database))

// catch 404 and forward to error handler
app.use(function(req, res, next) {
  let err = new Error('Not Found')
  err.status = 404
  next(err)
})

// error handler
app.use(function(err, req, res, next) {
  const status = err.status || 500

  res.status(status)
  res.render('error', {
    message: http.STATUS_CODES[status],
    error: req.app.get('env') === 'development' ? err : {}
  })
})

module.exports = app
