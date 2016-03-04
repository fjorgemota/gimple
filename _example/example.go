package main

import (
	"fmt"
	"net/http"

	"github.com/fjorgemota/gimple"
)

/** Lets create some interfaces so we can replace our APIs easily */
type AppInterface interface {
	Render(id uint64) string
}

type NameInterface interface {
	Name(id uint64) string
}

type PlanetInterface interface {
	Planet(id uint64) string
}

// Our Awesome Application renders a message using two APIs in our fake
// world.
type HomePlanetRenderApp struct {
	NameAPI   NameInterface
	PlanetAPI PlanetInterface
}

func (a *HomePlanetRenderApp) Render(id uint64) string {
	return fmt.Sprintf(
		"%s is from the planet %s.",
		a.NameAPI.Name(id),
		a.PlanetAPI.Planet(id),
	)
}

// Our fake Name API.
type NameAPI struct {
	HTTPTransport http.RoundTripper
}

func (n *NameAPI) Name(id uint64) string {
	// in the real world we would use f.HTTPTransport and fetch the name
	return "Spock"
}

// Our fake Planet API.
type PlanetAPI struct {
	HTTPTransport http.RoundTripper
}

func (p *PlanetAPI) Planet(id uint64) string {
	// in the real world we would use f.HTTPTransport and fetch the planet
	return "Vulcan"
}

type PlanetName struct {
	name string
}

func (p *PlanetName) Planet(id uint64) string {
	return p.name
}

// An example of provider just for planets

type ConfigPlanet struct{}

func (c ConfigPlanet) Register(g gimple.GimpleContainer) {
	// Here we'll define **just** planets
	g.Set("api.planet", func(container gimple.GimpleContainer) interface{} {
		transport := container.Get("transport").(http.RoundTripper)
		return &PlanetAPI{
			HTTPTransport: transport}
	})
	g.Set("api.planet.name", func(container gimple.GimpleContainer) interface{} {
		// Here we'll load an parameter
		return &PlanetName{name: container.Get("planet.name").(string)}
	})
}

// An example of provider
type Config struct{}

func (c Config) Register(g gimple.GimpleContainer) {
	/* Let's define some services! */
	g.Set("transport", func(container gimple.GimpleContainer) interface{} {
		return http.DefaultTransport
	})
	g.Set("planet.name", "Earth") // An example of parameter
	g.Set("api.name", func(container gimple.GimpleContainer) interface{} {
		transport := container.Get("transport").(http.RoundTripper)
		return &NameAPI{
			HTTPTransport: transport}
	})
	g.Set("app", func(container gimple.GimpleContainer) interface{} {
		// Try to replace "api.planet" with "api.planet.name"
		planet := container.Get("api.planet").(PlanetInterface)
		name := container.Get("api.name").(NameInterface)
		return &HomePlanetRenderApp{
			NameAPI:   name,
			PlanetAPI: planet}
	})
	// Load our planet configuration, note that we can load it AFTER defining our app
	g.Register(ConfigPlanet{})
}
func main() {
	// Create an instance of our Gimple Dependency Injection
	g := gimple.NewGimple()

	// In one line, we will setup our container provider :)
	g.Register(Config{})

	// Get our app and convert it to our interface
	app := g.Get("app").(AppInterface)

	fmt.Println(app.Render(42))
}
