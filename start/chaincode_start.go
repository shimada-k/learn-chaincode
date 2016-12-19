/*
Copyright IBM Corp 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"errors"
	"fmt"
    "encoding/json"
    "strings"
    "strconv"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type Family struct{
    Sex int `json:"sex"`
    Birthday string `json:"birthday"`
    SpouseId string `json:"spouse_id"`
    FatherId string `json:"father_id"`
    MotherId string `json:"monther_id"`
    ChildId string `json:"child_id"`
}

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

//==============
// Main
//==============
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	return nil, nil
}

func (t *SimpleChaincode) init_human(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    var err error
    id := args[0]
    sex, _ := strconv.Atoi(args[1])
    birthday := strings.ToLower(args[2])
    spouse_id := strings.ToLower(args[3])
    father_id := strings.ToLower(args[4])
    mother_id := strings.ToLower(args[5])
    child_id := strings.ToLower(args[6])

    str := `{"sex": ` + strconv.Itoa(sex) + `,"birthday": "` + birthday + `","spouse_id": "` + spouse_id + `","father_id": "` + father_id + `","mother_id": "` + mother_id + `","child_id": "` + child_id + `"}`

	fmt.Println(str)

    err = stub.PutState(id, []byte(str))
    if err != nil {
        return nil, err
    }
    return nil, nil
}

// Invoke is our entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "add" {     // 人を追加する
		return t.init_human(stub, args)
	} else if function == "hospital" {   // 親の子供IDを消す
        familyAsBytes, err := stub.GetState(args[0])
        if err != nil {
            return nil, err
        }

        res := Family{}
        json.Unmarshal(familyAsBytes, &res)
        res.ChildId = ""    // 子供を消す

        jsonAsBytes, _ := json.Marshal(res)
        stub.PutState(args[0], jsonAsBytes)
        return nil, nil
    } else if function == "pluged" {   // 子供の親IDを消す
        familyAsBytes, err := stub.GetState(args[0])
        if err != nil {
            return nil, err
        }

        res := Family{}
        json.Unmarshal(familyAsBytes, &res)
        res.FatherId = ""    // 父親を消す
        res.MotherId = ""    // 母親を消す

        jsonAsBytes, _ := json.Marshal(res)
        stub.PutState(args[0], jsonAsBytes)
        return nil, nil
    } else if function == "adopted" {   // 子供の親IDを上書きする
        familyAsBytes, err := stub.GetState(args[0])
        if err != nil {
            return nil, err
        }

        res := Family{}
        json.Unmarshal(familyAsBytes, &res)
        res.FatherId = args[1]    // 父親IDを上書きする

        jsonAsBytes, _ := json.Marshal(res)
        stub.PutState(args[0], jsonAsBytes)
        return nil, nil
    }
	fmt.Println("invoke did not find func: " + function)    //error

	return nil, errors.New("Received unknown function invocation: " + function)
}

func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    var key, jsonResp string
    var err error

    if len(args) != 1 {
        return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
    }

    key = args[0]
    valAsbytes, err := stub.GetState(key)
    if err != nil {
        jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
        return nil, errors.New(jsonResp)
    }

    return valAsbytes, nil
}

func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
    fmt.Println("query is running " + function)

    // Handle different functions
    if function == "read" {                            //read a variable
        return t.read(stub, args)
    }
    fmt.Println("query did not find func: " + function)

    return nil, errors.New("Received unknown function query: " + function)
}

