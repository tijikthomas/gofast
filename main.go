package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type Person struct {
	Name         string         `json:"name"`
	Age          int            `json:"age"`
	IsLearningGo bool           `json:"is_learning_go"`
	Skills       map[string]int `json:"skills"`
}

var people []Person
var mutx sync.RWMutex

func main() {
	http.HandleFunc("/hello", handleHello)
	http.HandleFunc("/people", getPeople)
	http.HandleFunc("/create", createPerson)
	http.HandleFunc("/people/{name}", handlePeople)
	//http.HandleFunc("/people/{name}", updatePerson)

	fmt.Println("Server is running on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}

}

func handleHello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, Go!"))
}

func getPeople(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	mutx.RLock()
	defer mutx.RUnlock()

	peopleCopy := make([]Person, len(people))
	copy(peopleCopy, people)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(peopleCopy)
}

func createPerson(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var person Person
	err := json.NewDecoder(r.Body).Decode(&person)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid input %v", err), http.StatusBadRequest)
		return
	}

	mutx.Lock()
	defer mutx.Unlock()

	people = append(people, person)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Person added successfully"})

}

func deletePerson(target_index int) {
	mutx.Lock()
	defer mutx.Unlock()

	people[target_index] = people[len(people)-1]
	people = people[:len(people)-1]
}

func updatePerson(target_index int, params map[string]interface{}) {
	fmt.Println(params)
	mutx.Lock()
	defer mutx.Unlock()

	person := people[target_index]
	for key, value := range params {
		switch key {
		case "name":
			if v, ok := value.(string); ok {
				person.Name = v
			}
		case "age":
			if v, ok := value.(int); ok {
				person.Age = v
				fmt.Println(person.Age)
			}
		case "isLearningGo":
			if v, ok := value.(bool); ok {
				person.IsLearningGo = v
			}
		case "skills":
			if v, ok := value.(map[string]interface{}); ok {
				newSkills := make(map[string]int)
				for skill, level := range v {
					if levelFloat, ok := level.(float64); ok {
						newSkills[skill] = int(levelFloat)
					}
				}
				person.Skills = newSkills

			}
		}
	}
	people[target_index] = person
	fmt.Println(target_index, people)
}

func handlePeople(w http.ResponseWriter, r *http.Request) {
	var name string
	var found bool
	name = r.PathValue("name")
	var target_index int
	for index, person := range people {
		if person.Name == name {
			found = true
			target_index = index
			break
		}
	}
	if found != true {
		http.Error(w, "Person not found", http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodPatch:
		var params map[string]interface{}
		json.NewDecoder(r.Body).Decode(&params)
		fmt.Println(params)
		updatePerson(target_index, params)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Person updated successfully"})

	case http.MethodDelete:
		deletePerson(target_index)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Person deleted successfully"})

	default:
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}

}
