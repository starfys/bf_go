//This is a brainfuck interpreter written in GO
//Copyright 2016 Steven Sheffey
/*
    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/
package main

import (
    "container/list"
	"errors"
    "io/ioutil"
	"os"
    "strconv"
)
func main() {

//Constants for the interpreter
	//The size of the memory space
	const MEMORY_SIZE int = 16384
//VM variables:
	//The data space
	var memory [MEMORY_SIZE]byte
	//The data pointer
	var dataPointer int = 0
	//The instruction pointer
	var instructionPointer int = 0
	//The loop stack
    loopStack := list.New()
//Interpreter variables
	//The brainfuck executable, read from a file
	var program []byte
	//Used to store errors
	var errorHandler error
	//Used for {get,put}char
	curChar := []byte{0}
	//Initialization
	//Ensure there command line includes enough space for a simple command line
	if len(os.Args) != 2 {
		panic(errors.New("Usage: interp filename.bf"))
	}
	//Read the program from the file
	program, errorHandler = ioutil.ReadFile(os.Args[1])

	//Check to see if there was an error in reading the file
	if errorHandler != nil {
		panic(errorHandler)
	}

    //Initialize the memory
	for dataPointer = 0; dataPointer < MEMORY_SIZE; dataPointer++ {
	    memory[dataPointer] = 0
    }
    //Reset dataPointer to 0
    dataPointer = 0
    //Develop a loop reference map
    loopMap := map[int]int{}
    for instructionPointer = 0; instructionPointer < len(program); instructionPointer++ {
        switch program[instructionPointer] {
            case '[':
                loopStack.PushBack(instructionPointer)
            case ']':
                if loopStack.Len() > 0 {
                    loopMap[loopStack.Back().Value.(int)] = instructionPointer
                    loopMap[instructionPointer] = loopStack.Back().Value.(int)
                    loopStack.Remove(loopStack.Back())
                } else {
                    panic(errors.New("Extra ']' at position " + strconv.Itoa(instructionPointer)))
                }
        }

    }
    if loopStack.Len() > 0 {
        panic(errors.New("Unmatched '[' at position " + strconv.Itoa(loopStack.Back().Value.(int))))
    }

    //Iterate through each token of the program
    for instructionPointer = 0; instructionPointer < len(program); instructionPointer++ {
        switch program[instructionPointer] {
            //Move current position forward
            case '>':
                dataPointer = dataPointer + 1
                if dataPointer == len(program) {
                    panic(errors.New("Error: Memory went out of bounds at position " + strconv.Itoa(instructionPointer)))
                }
            //Move current position backward
            case '<':
                dataPointer = dataPointer - 1
                if dataPointer == -1 {
                    panic(errors.New("Error: Memory went out of bounds at position " + strconv.Itoa(instructionPointer)))
                }
            //Increment memory at at the current position in memory
            case '+':
                memory[dataPointer] = memory[dataPointer] + 1
            //Decrement memory at the current position in memory
            case '-':
                memory[dataPointer] = memory[dataPointer] - 1
            case ',':
                os.Stdin.Read(curChar)
                memory[dataPointer] = curChar[0]
            case '.':
                curChar[0] = memory[dataPointer]
                os.Stdout.Write(curChar)
            case '[':
                if memory[dataPointer] == 0 {
                    instructionPointer = loopMap[instructionPointer]
                }
            case ']':
                if memory[dataPointer] != 0 {
                    instructionPointer = loopMap[instructionPointer]
                }
        }
    
    }

}
