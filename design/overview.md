# The idea

There is too much boiler plate. It would be great to be able to bootstrap a project in an afternoon and write only the business logic. I have a great idea, a chat application for business users let's call it slack. What do I need to build.

Frameworks like Django all require the same code be rewritten to setup your application. Further more support of newer protocols such as widely used Websockets and the fast approaching HTTP 2 are not supported out of the box.

So lets have the following features,

- Locally runnable, you can see and modify the code as you wish

- Implement an MVC pattern, Framework user only implements the models, views and controllers. They should already be plumbed together.

  ![img](https://upload.wikimedia.org/wikipedia/commons/thumb/a/a0/MVC-Process.svg/500px-MVC-Process.svg.png)

  - Models are designed in a simple markup, JSON or YAML something like that
    - Migrations are done for you, you should be able to inject SQL if need be.
  - Controllers are presented with standard interfaces which can be programmed
    - You should be able to choose your implementation langauge to implement this but this is a definite stretch feature. It will all be go for the moment. If this is going to be useful we can implement this.
  - Views are auto generated as a CRUD API on the models.  These auto generated views must be able to be edited.
    - How are we going to route?
    - Connections can be upgraded to a websocket by adding a flag.
    - A JS SDK should be made available to make it easier.
    - Roll based access control should be implemented

- Live notifications of events in the system

  - Web socket support to hook into your front end

- Inject standard and customisable middleware

  - It is useful, but do it later.



Things to consider:

- CORS
- Sessions
- CSRF

Not very insightful, needs to be fleshed out



## Useful packages and inspiration

- Auth provider - https://github.com/markbates/goth

- Wrest - https://www.npmjs.com/package/wrest




