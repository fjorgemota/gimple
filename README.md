# Gimple

[![Build Status](https://travis-ci.org/fjorgemota/gimple.svg)](https://travis-ci.org/fjorgemota/gimple)

This project is a "port" of [Pimple Dependency Injection container](https://github.com/silexphp/Pimple/) to Go.

The code of this project may not be in his state-of-art, but seems to be a great start to learn a few more about Go, and to avoid initializing a bunch of things in your application manually.

All the code is tested using **Ginkgo** and seems to be stable. Below is the documentation for the project:

## Features

Good projects have good features. And because this here's the list of features that Gimple supports:

- Define services;
- Define factories;
- Define parameters easily;
- Defining services/parameters/factories from another files - because you should be able to split your configuration easily;
- Simple API;
- Allows extending services easily;
- Allow to get the raw service creator easily;
- Pure Go, no C code envolved;
- [Fully tested](https://travis-ci.org/fjorgemota/gimple) on each commit;
- I already said that it have a really Simple API? :)

## Installation

The installation of this package is very simple: In fact, it can be installed by just running:

```
    go get github.com/fjorgemota/gimple
```

## Usage

Creating a Gimple Container is just a matter of creating a Gimple instance:

```go
import "gimple"
container := gimple.NewGimple()
```

Gimple, as Pimple and many other dependency injections containers, manage two different kind of data: **services** and **parameters**.

## Defining services

As Pimple describes, a service is an object that does something as part of a larger system. Examples of services: a database connection, a templating engine, or a mailer. Almost any global object can be a service.

Services in Gimple (and in Pimple too!) are defined by anonymous functions that return an instance of an object. Different from Pimple, however, here we need to call the method `Set()` on Gimple container, as Go offers no way to simulate map like behavior as in Pimple:

```go
// define some services
container.Set('session_storage', func (c gimple.GimpleContainer) {
    return SessionStorage{'SESSION_ID'};
});

container.Set('session', func (c gimple.GimpleContainer) {
    session_storage := c.Get('session_storage').(SessionStorage)
    return Session{};
});
```

Notice that the anonymous function that define a service has access to the current container instance, allowing references to other services or parameters.

The objects are created on demand, just when you get them. The order of the definitions does not matter.

Using the defined services is very easy, too:

```go
// get the session object
session := container.Get('session').(Session);

// the above call is roughly equivalent to the following code:
// storage = SessionStorage{'SESSION_ID'};
// session = Session{storage};
```

## Defining factory services

By default, when you get a service, Gimple automatically cache it's value, returning always the **same instance** of it. If you want a different instance to be returned for all calls, wrap your anonymous function with the `Factory()` method:

```go
container.Set('session', container.Factory(func (c gimple.GimpleContainer) interface{} {
    session_storage := c.Get('session_storage').(SessionStorage)
    return Session{session_storage};
}));
```

Now, each time you call `container.Get('session')`, a new instance of `Session` is returned for you.

## Defining parameters

Defining a parameter allows to ease the configuration of your container from the outside and to store global values. In Gimple, parameters are defined as anything that it's not a function:

```go
// define a parameter called cookie_name
container.Set('cookie_name', 'SESSION_ID');
```

If you change the `session_storage` service definition like below:

```go
container.Set('session_storage', func (c gimple.GimpleContainer) {
    cookie_name := c.Get('cookie_name').(string)
    return SessionStorage{cookie_name};
});
```

You can now easily change the cookie name by overriding the `cookie_name` parameter instead of redefining the service definition.

### Defining parameters based on environment variables

Do you wanna do define parameters in the container based on environment variables? It's okay! You can define it easily like that:

```go
import "os" // At the top of the file
//define parameter based on environment variable
container.Set('cookie_name', os.Getenv(COOKIE_NAME));
```

## Protecting parameters

Because Gimple see anything that is a function as a service, you need to wrap anonymous functions with the `Protect()` method to store them as parameters:

```go
import "math/rand" // At the top of the file
container.Set('random_func', container.Protect(func () {
    return rand.Int();
}));
```

## Modifying Services after Definition

In some cases you may want to modify a service definition after it has been defined. You can use the `Extend()` method to define additional code to be run on your service just after it is created:

```go
container.Set('session_storage', func (c gimple.GimpleContainer) {
    cookie_name := c.Get('cookie_name').(string)
    return SessionStorage{cookie_name};
});

container.Extend('session_storage', func (old interface{}, c gimple.GimpleContainer) interface{} {
    storage := old.(SessionStorage)
    storage.SomeMethod();

    return storage;
});
```

The first argument is the name of the service to extend, the second a function that gets access to the object instance and the container.

## Extending a Container

If you use the same libraries over and over, you might want to reuse some services from one project to the next one; package your services into a provider by implementing the following object structure by duck-typing:

```go
type provider struct{}

func (p *provider) Register(c gimple.GimpleContainer) {
	// Define your services and parameters here
}
```

After creating a object with that structure, you can register it in the container:

```go
container.Register(provider{});
```

## Fetching the Service Creation Function

When you access an object, Gimple automatically calls the anonymous function that you defined, which creates the service object for you. If you want to get raw access to this function, but don't want to `protect()` that service, you can use the `raw()` method to access the function directly:

```go
container.Set('session', func (c gimple.GimpleContainer) interface{} {
    storage := c.Get('session_storage').(SessionStorage)
    return Session{storage};
});

sessionFunction := container.Raw('session').(Session);
```

## Last, but not least important: Customization

Do you wanna to customize Gimple's functionally? You can! Just extend it using ES6's class syntax:

```go
var Gimple = require("Gimple");

type ABigContainer struct{
    *Gimple
}

// Overwrite any of the Gimple's methods here

var container = ABigContainer}{}; 
```

Good customization. :)