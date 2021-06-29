package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type SmartContract struct{}
type User struct {
	ID                string         `json:"id"`
	Avg               float32        `json:"avg"`
	CurrentProjects   []string       `json:"current_projects"`
	CompletedProjects []string       `json:"completed_projects"`
	ProjScoreMap      map[string]int `json:"proj_score_map"`
}
type Project struct {
	ID string `json:"id"`
	// Name         string   `json:"name"`
	Participants []string `json:"participants"`
	Completed    bool     `json:"completed"`
}

type resultFormat struct {
	success bool
	result  []byte
	error   string
}

func (s *SmartContract) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}
func (s *SmartContract) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	fn, args := stub.GetFunctionAndParameters()
	if fn == "registerUser" {
		return s.registerUser(stub, args)
	} else if fn == "registerProject" {
		return s.registerProject(stub, args)
	} else if fn == "joinProject" {
		return s.joinProject(stub, args)
	} else if fn == "completeProject" {
		return s.completeProject(stub, args)
	} else if fn == "recordScore" {
		return s.recordScore(stub, args)
	} else if fn == "getUserInfo" {
		return s.getUserInfo(stub, args)
	} else {
		return shim.Error("Not Supported Function name")
	}
}

func (s *SmartContract) registerUser(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 { // usr_id
		return shim.Error("registerUser uses 1 parameter")
	}
	// error if id exist
	usr_id := args[0]
	result, err := stub.GetState(usr_id)
	if err == nil && result == nil {
		usr := User{ID: usr_id, Avg: 0., CurrentProjects: make([]string, 0), CompletedProjects: make([]string, 0), ProjScoreMap: make(map[string]int)}
		usrAsBytes, _ := json.Marshal(usr)
		stub.PutState(usr_id, usrAsBytes)
		return shim.Success([]byte("Register User"))
	} else {
		return shim.Error("user id is duplicated")
	}

}

func (s *SmartContract) registerProject(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2 { // proj_id, usr_id
		return shim.Error("registerProject uses 2 parameter")
	}
	proj_id := args[0]
	usr_id := args[1]
	// error if id exist
	check_proj, perr := stub.GetState(proj_id)
	check_usr, uerr := stub.GetState(usr_id)
	if perr == nil && uerr == nil && check_proj == nil && check_usr != nil {
		p := make([]string, 1)
		p[0] = usr_id
		proj := Project{ID: proj_id, Participants: p, Completed: false}
		projAsBytes, _ := json.Marshal(proj)
		stub.PutState(proj_id, projAsBytes)
		return shim.Success([]byte("Register Project"))
	} else {
		return shim.Error("proj id(" + proj_id + ") is duplicated or there is no user(" + usr_id + ")")
	}

}

func (s *SmartContract) joinProject(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2 { // proj_id, usr_id
		return shim.Error("joinProject uses 2 parameter")
	}
	proj_id := args[0]
	usr_id := args[1]
	projAsBytes, perr := stub.GetState(proj_id)
	usrAsBytes, uerr := stub.GetState(usr_id)
	if perr == nil && uerr == nil && projAsBytes != nil && usrAsBytes != nil {
		proj := Project{}
		_ = json.Unmarshal(projAsBytes, &proj)
		usr := User{}
		_ = json.Unmarshal(usrAsBytes, &usr)
		proj.Participants = append(proj.Participants, usr_id)
		projAsBytes, _ = json.Marshal(proj)
		usr.CurrentProjects = append(usr.CurrentProjects, proj_id)
		usrAsBytes, _ = json.Marshal(usr)

		stub.PutState(proj_id, projAsBytes)
		stub.PutState(usr_id, usrAsBytes)

		return shim.Success([]byte("Join project"))
	} else {
		return shim.Error("There is error proj id(" + proj_id + ") or user(" + usr_id + ")")
	}
}

func (s *SmartContract) completeProject(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 { // proj_id
		return shim.Error("joinProject uses 2 parameter")
	}
	proj_id := args[0]
	projAsBytes, err := stub.GetState(proj_id)
	if err == nil && projAsBytes != nil {
		proj := Project{}
		_ = json.Unmarshal(projAsBytes, &proj)
		for i := 0; i < len(proj.Participants); i++ {
			usrAsBytes, _ := stub.GetState(proj.Participants[i])
			usr := User{}
			_ = json.Unmarshal(usrAsBytes, &usr)
			usr.CurrentProjects = deleteStringInSlice(proj_id, usr.CurrentProjects)
			usr.CompletedProjects = append(usr.CompletedProjects, proj_id)

			usrAsBytes, _ = json.Marshal(usr)
			stub.PutState(proj.Participants[i], usrAsBytes)
		}
		proj.Completed = true
		projAsBytes, _ = json.Marshal(proj)
		stub.PutState(proj_id, projAsBytes)
		return shim.Success([]byte("Complete project"))
	}
	return shim.Error("There is error proj id(" + proj_id + ")")
}

func (s *SmartContract) recordScore(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 3 { // proj_id, usr_id, score
		return shim.Error("recordScore uses 3 parameter")
	}
	proj_id := args[0]
	usr_id := args[1]
	projAsBytes, perr := stub.GetState(proj_id)
	usrAsBytes, uerr := stub.GetState(usr_id)
	if uerr != nil || perr != nil {
		return shim.Error("GetState Error")
	}
	if usrAsBytes == nil || projAsBytes == nil {
		return shim.Error("No data")
	}
	proj := Project{}
	json.Unmarshal(projAsBytes, &proj)
	if proj.Completed {
		score, err := strconv.Atoi(args[2])
		if err != nil {
			return shim.Error("Score is not int")
		}
		usr := User{}
		_ = json.Unmarshal(usrAsBytes, &usr)
		usr.ProjScoreMap[proj_id] = score
		proj_num := len(usr.CompletedProjects)
		usr.Avg = (usr.Avg*float32(proj_num) + float32(score)) / float32(proj_num+1)

		usrAsBytes, _ = json.Marshal(usr)
		stub.PutState(usr_id, usrAsBytes)

		return shim.Success([]byte("Score is recorded"))
	} else {
		return shim.Error("Project is Not completed")
	}
}

func deleteStringInSlice(a string, arr []string) []string {
	if stringInSlice(a, arr) {
		for i, v := range arr {
			if v == a {
				arr = remove(arr, i)
				return arr
			}
		}
	}
	return arr
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
func remove(s []string, i int) []string {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func (s *SmartContract) getUserInfo(APIstub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 { //user_id
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	UserAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(UserAsBytes)
}

func main() {
	err := shim.Start((new(SmartContract)))
	if err != nil {
		fmt.Printf("Error starting chaincode : %s", err)
	}
}
