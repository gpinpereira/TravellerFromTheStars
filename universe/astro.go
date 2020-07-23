package universe

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"time"

	pb "../services"
)

//import "fmt"

type Astro struct {
	Name                             string
	Age                              int
	is_sun                           bool
	is_ship                          bool
	x, y, vx, vy, mass, radius       float64
	next_x, next_y, next_vx, next_vy float64
	center_body                      string
	system                           SolarSystem
}

type SolarSystem struct {
	Name string
	//bodies []*Astro
	Age    float64
	deltat float64
	bodies map[string]*Astro
}

type SolarSystemInterface interface {
	GetBodyByName()
}

func (sys *SolarSystem) GetBodies() map[string]*Astro {

	return sys.bodies
}

func (sys *SolarSystem) GetAge() float64 {
	//fmt.Println("Age")
	//fmt.Println(sys.Age)
	return sys.Age
}

type AstroNextPosition interface {
	rk4()
	nextStep()
}

func (astro *Astro) GetStats() (float64, float64, float64, float64, float64) {
	return astro.x, astro.y, astro.vx, astro.vy, astro.mass

}

func (astro *Astro) calc_acel(x float64, y float64) (float64, float64) {
	otherbodies := astro.system.GetBodies()

	//body := otherbodies[astro.center_body]
	ax := 0.0
	ay := 0.0

	var G float64 = 6.67408 * math.Pow10(-11)

	///fmt.Println("!!!!!!!!!!!!!!!!!!!!")
	//fmt.Println(astro.Name)
	//fmt.Println(body.Name)
	//fmt.Println(body.Name)

	for _, body := range otherbodies {
		if body.Name != astro.Name {
			//if astro.Name == "Luna" {
			//	fmt.Println(body.Name)
			//}

			var dx float64 = x - body.x
			var dy float64 = body.y - astro.y
			var dsq float64 = dx*dx + dy*dy
			var dr float64 = math.Sqrt(dsq)

			var force float64 = G * body.mass / dsq
			//fmt.Println("forcex", force)
			ax += force * dx / dr

			//for ay

			dx = body.x - astro.x
			dy = y - body.y
			dsq = dx*dx + dy*dy
			dr = math.Sqrt(dsq)

			force = G * body.mass / dsq
			ay += force * dy / dr

		}
	}
	return -ax, -ay
}

func (astro *Astro) rk4() {

	var deltat float64 = astro.system.deltat
	//fmt.Println(astro.Name)

	k1x, k1y := astro.calc_acel(astro.x, astro.y)
	k1x = k1x * deltat
	k1y = k1y * deltat

	l1x := deltat * astro.vx
	l1y := deltat * astro.vy

	k2x, k2y := astro.calc_acel(astro.x+0.5*l1x, astro.y+0.5*l1y)
	k2x = k2x * deltat
	k2y = k2y * deltat

	l2x := deltat * (astro.vx + 0.5*k1x)
	l2y := deltat * (astro.vy + 0.5*k1y)

	k3x, k3y := astro.calc_acel(astro.x+0.5*l2x, astro.y+0.5*l2y)
	k3x = k3x * deltat
	k3y = k3y * deltat

	l3x := deltat * (astro.vx + 0.5*k2x)
	l3y := deltat * (astro.vy + 0.5*k3y)

	k4x, k4y := astro.calc_acel(astro.x+l3x, astro.y+l3y)
	k4x = k4x * deltat
	k4y = k4y * deltat

	l4x := deltat * (astro.vx + k3x)
	l4y := deltat * (astro.vy + k3y)

	//fmt.Println(l1x, l2x, l3x, l4x)
	//fmt.Println(l1y, l2y, l3y, l4y)
	astro.next_x = astro.x + (l1x+2*l2x+2*l3x+l4x)/6
	astro.next_y = astro.y + (l1y+2*l2y+2*l3y+l4y)/6

	astro.next_vx = astro.vx + (k1x+2*k2x+2*k3x+k4x)/6
	astro.next_vy = astro.vy + (k1y+2*k2y+2*k3y+k4y)/6
}

func (astro *Astro) nextStep() {

	astro.x = astro.next_x
	astro.y = astro.next_y

	astro.vx = astro.next_vx
	astro.vy = astro.next_vy
}

func moveAstro(nextp AstroNextPosition) {

	nextp.rk4()
}

func takeNextStep(nextp AstroNextPosition) {

	nextp.nextStep()

}

func MakeSystemCSV() *SolarSystem {

	csvfile, err := os.Open("solar_pos.csv")
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}

	// Parse the file
	r := csv.NewReader(csvfile)
	r.Comma = ';'
	//dt=1
	s := SolarSystem{Name: "Solar System", Age: 0, deltat: 1}
	bodies := make(map[string]*Astro)

	i := 0
	for {
		// Read each record from csv
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		name := record[0]
		if name == "" || i == 0 {
			i += 1
			continue
		}
		x := StringToFloat(record[1]) * 1000
		y := StringToFloat(record[2]) * 1000
		//z := StringToFloat(record[3])
		vx := StringToFloat(record[4]) * 1000
		vy := StringToFloat(record[5]) * 1000
		//vz := StringToFloat(record[6])
		mass := StringToFloat(record[7])
		center_body := record[8]
		//fmt.Println(name, x, y, vx, vy, mass)
		if name == "Sun" {
			bod := Astro{Name: name, Age: 0, x: x, y: y, vx: vx, vy: vy, mass: mass, is_sun: true, is_ship: false}
			bod.system = s
			bodies[name] = &bod

		} else {
			bod := Astro{Name: name, Age: 0, x: x, y: y, vx: vx, vy: vy, mass: mass, is_sun: false, is_ship: false, center_body: center_body}
			bod.system = s
			bodies[name] = &bod

		}

		if name == "Earth" {

			bod := Astro{Name: "Ship", Age: 0, x: x * 0.99, y: y * 0.99, vx: vx, vy: vy, mass: 200, is_sun: false, is_ship: true, center_body: "Sun"}
			bod.system = s
			bodies["Ship"] = &bod

		}

	}

	s.bodies = bodies

	for i := range s.bodies {

		fmt.Println(s.bodies[i].Name)
		s.bodies[i].system = s
	}

	return &s
}

func MakeSystem() SolarSystem {

	star := Astro{Name: "Sun", Age: 0, x: 0, y: 0, vx: 0, vy: 0, mass: 1.989 * math.Pow10(30), is_sun: true}
	//planet1 := Astro{Name: "Earth", Age: 0, x: 4.48456 * math.Pow10(10), y: 1.40453 * math.Pow10(11), vx: -28862.6, vy: 8959.90, mass: 5.97 * math.Pow10(24), is_sun: false}
	planet1 := Astro{Name: "Earth", Age: 0, x: -2.81376 * math.Pow10(9), y: 1.47128 * math.Pow10(11), vx: -30274.75, vy: -689.318, mass: 5.97 * math.Pow10(24), is_sun: false}
	planet2 := Astro{Name: "Venus", Age: 0, x: -1.63282 * math.Pow10(10), y: -1.07443 * math.Pow10(11), vx: -34386.74, vy: 5400.21, mass: 4.86 * math.Pow10(24), is_sun: false}

	planet3 := Astro{Name: "Mars", Age: 0, x: -2.414562 * math.Pow10(11), y: -4.1365 * math.Pow10(10), vx: 4995.903, vy: -21811.701, mass: 6.42 * math.Pow10(23), is_sun: false}

	planet4 := Astro{Name: "Jupiter", Age: 0, x: -6.436111 * math.Pow10(11), y: -4.962076837 * math.Pow10(11), vx: 7828.37152, vy: -9740.141258, mass: 1.89727 * math.Pow10(27), is_sun: false}

	bodies := make(map[string]*Astro)
	bodies["Sun"] = &star
	bodies["Earth"] = &planet1
	bodies["Venus"] = &planet2
	bodies["Mars"] = &planet3
	bodies["Jupiter"] = &planet4
	s := SolarSystem{Name: "Solar System", Age: 0, deltat: 0}
	s.bodies = bodies
	star.system = s
	planet1.system = s
	planet2.system = s
	planet3.system = s
	planet4.system = s
	return s
}

func SimulateSystem(system *SolarSystem) {

	//system := make_system()
	time_tick := 0
	this_system_bodies := system.bodies

	for true {
		//if start
		for _, body := range this_system_bodies {

			if !body.is_sun {
				//fmt.Println(body)
				moveAstro(body)
			}
		}

		for _, body := range this_system_bodies {

			if !body.is_sun {
				//fmt.Println(body)
				takeNextStep(body)
			}
		}

		system.Age = system.Age + system.deltat
		time_tick += 1

		time.Sleep(2 * 6900 * time.Nanosecond) // for dt =1
		//time.Sleep(2 * time.Nanosecond)

	}

}

func FastSimulation(in *pb.MyState) ([]float64, []float64) {

	var x_hist []float64
	var y_hist []float64

	allbodies := in.Otherbodies.Bodies

	dt := 1800.0
	system := SolarSystem{Name: "Solar System", Age: 0, deltat: dt}
	bodies := make(map[string]*Astro)

	fmt.Println("starting fast simulation")

	for name, body := range allbodies {

		the_sun := false
		if name == "Sun" {
			the_sun = true
		}
		bod := Astro{Name: name, Age: 0, x: body.X, y: body.Y, vx: body.Vx, vy: body.Vy, mass: body.Mass, is_sun: the_sun, is_ship: false, center_body: "Sun"}
		bod.system = system
		bodies[name] = &bod
		//send_bodies[name] = &pb.BodyPos{X:, Y: body.X, Vx: body.Vx, Vy: body.Vy, Mass: body.mass}

	}

	system.bodies = bodies

	this_system_bodies := system.bodies
	limit := 3600.0 * 24 * 365

	t := 0.0
	fmt.Println(t)
	for t < limit {
		//if start
		fmt.Println(t)
		for _, body := range this_system_bodies {

			if !body.is_sun {
				//fmt.Println(body)
				moveAstro(body)
			}
		}

		for _, body := range this_system_bodies {

			if !body.is_sun {
				//fmt.Println(body)
				takeNextStep(body)
			}
		}

		system.Age = system.Age + system.deltat
		t += dt

	}

	fmt.Println(t)
	fmt.Println("end fast simulation")

	return x_hist, y_hist

}

func StringToFloat(f string) float64 {
	s, err := strconv.ParseFloat(f, 64)

	if err == nil {
		return s
	} else {

		fmt.Println(err)
	}
	return -1.0

}
