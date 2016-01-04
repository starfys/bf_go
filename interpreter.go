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
	"errors"
	"io/ioutil"
	"os"
    "fmt"
)

func main() {
	//Constants for the interpreter
	//TODO: change these through command line flags
	//The size of the memory space
	const MEMORY_SIZE int = 1024
	//The maximum depth of [] loops
	const MAX_DEPTH int = 1024
	//VM variables:
	//The data space
	//TODO: allow the option of a dynamic array
	//TODO: configure this through command line flags
	var memory [MEMORY_SIZE]byte
	//The data pointer
	var dataPointer int = 0
	//The instruction pointer
	var instructionPointer int = 0
	//The loop stack
	//TODO: set this up to be a real stack
	var loopStack [MAX_DEPTH]int
	//Position of the last loop started
	var lastLoop int = -1
    //Needed for testing nests
    var bracketBalance int
	//Interpreter variables
	//The brainfuck executable, read from a file
	var program []byte
	//Used to store errors
	var errorHandler error
	//Used for {get,put}char
	curChar := []byte{0}
	//Used for iterating arrays
	var index int
	//Initialization
	//Ensure there command line includes enough space for a simple command line
	if len(os.Args) != 2 {
		panic(errors.New("Usage: interp filename.bf"))
	}
	//Read the program from the file
	program, errorHandler = ioutil.ReadFile(os.Args[1])
	for index = 0; index < MEMORY_SIZE; index++ {
		memory[index] = 0
	}
	//Check to see if there was an error in reading the file
	if errorHandler != nil {
		panic(errorHandler)
	}

	//Iterate through the brainfuck code
	for instructionPointer = 0; instructionPointer < len(program); instructionPointer++ {
        fmt.Println(instructionPointer)
		//Switch for each token
		switch program[instructionPointer] {
		//Move data pointer forward
		case '>':
			//Check for overflow
			if dataPointer == MEMORY_SIZE {
				panic(errors.New("Data pointer overflow."))
			}
			dataPointer = dataPointer + 1

		//Move data pointer backward
		case '<':
			//Check for underflow
			if dataPointer == -1 {
				panic(errors.New("Data pointer underflow."))
			}
			dataPointer = dataPointer - 1

		//Increment memory at current data pointer
		case '+':
			memory[dataPointer] += 1

		//Decrement memory at current data pointer
		case '-':
			memory[dataPointer] -= 1

		//Read character from stdin
		case ',':
			os.Stdin.Read(curChar)
			memory[dataPointer] = curChar[0]

		//Write character to stdout
		case '.':
			curChar[0] = memory[dataPointer]
			os.Stdout.Write(curChar)

		//Jumps to the next ']' if the data pointer is zero
		case '[':
			//TODO: check for mismatched brackets
			//End loop
			if memory[dataPointer] == 0 {
				//Iterate instruction pointer until the opposite ']' is reached
				for bracketBalance = 1; bracketBalance > 0; instructionPointer++ {
                    //TODO: mismatched bracket errors here
                    if program[instructionPointer] == '[' {
                        bracketBalance = bracketBalance + 1
                    }
                    if program[instructionPointer] == ']' {
                        bracketBalance = bracketBalance - 1
                    }
                }
			} else { //Continue with loop
				//Store the location of the loop start
				lastLoop = lastLoop + 1
				loopStack[lastLoop] = instructionPointer
			}

		//End the loop if the data at the current pointer is zero
		case ']':
			//TODO: check for mismatched brackers
			//End loop
			if memory[dataPointer] == 0 {
				loopStack[lastLoop] = 0
				lastLoop = lastLoop - 1
			} else { //Restart loop
				instructionPointer = loopStack[lastLoop]
			}

		}
	}
}
